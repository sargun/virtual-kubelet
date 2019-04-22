package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/virtual-kubelet/virtual-kubelet/providers"
	"github.com/virtual-kubelet/virtual-kubelet/providers/plugin/proto"
	"github.com/virtual-kubelet/virtual-kubelet/providers/plugin/shared"
	"os"
	"os/exec"
	"regexp"

	"github.com/hashicorp/go-plugin"
)

// This is just a restrictive to avoid any weirdness in finding the plugin name in the path
var pluginNameRegex = regexp.MustCompile("^[[:lower:]][a-z0-9-_]*$")

// Configuration is the configuration for the plugin provider
type Configuration struct {
	// PluginName is the name of the plugin.
	// The plugin binary is expected to be in the named: virtual-kubelet-plugin-NAME
	PluginName string `json:"pluginName,omitempty"`
}

type simpleProvider struct {
	providers.Provider
}

type providerWithPodNotifier struct {
	providers.Provider
	providers.PodNotifier
}

// NewPluginProvider creates a new PluginProvider
func NewPluginProvider(providerConfig, nodeName, operatingSystem string, internalIP string, daemonEndpointPort int32) (providers.Provider, error) {
	var podNotifierProvider providers.PodNotifier

	var config Configuration
	if err := loadConfig(providerConfig, &config); err != nil {
		return nil, err
	}



	executableName := fmt.Sprintf("virtual-kubelet-plugin-%s", config.PluginName)
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.HandshakeConfig(config.PluginName),
		VersionedPlugins: map[int]plugin.PluginSet{
			shared.VersionWithFeatures(1): map[string]plugin.Plugin{
				shared.ProviderPluginName:  &providerPlugin{},
			},
			shared.VersionWithFeatures(1, shared.PodNotifier): map[string]plugin.Plugin{
				shared.ProviderPluginName:    &providerPlugin{},
				shared.PodNotifierPluginName: &podNotifierPlugin{},
			},
		},
		Cmd: exec.Command(executableName),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Managed: true,
	})


	protocolClient, err := pluginClient.Client()
	if err != nil {
		return nil, err
	}

	rawProviderPluginClient, err := protocolClient.Dispense(shared.ProviderPluginName)
	if err != nil {
		return nil, err
	}

	providerPlugin := rawProviderPluginClient.(*ProviderPlugin)
	_, err = providerPlugin.client.Register(context.TODO(), &proto.ProviderRegisterRequest{
		InitConfig: &proto.InitConfig{
			ConfigPath:      providerConfig,
			NodeName:        nodeName,
			OperatingSystem: operatingSystem,
			InternalIP:      internalIP,
			DaemonPort:      daemonEndpointPort,
		},
	})

	if shared.HasFeature(pluginClient.NegotiatedVersion(), shared.PodNotifier) {
		rawProviderPodNotifierClient, err := protocolClient.Dispense(shared.PodNotifierPluginName)
		if err != nil {
			return nil, err
		}

		podNotifierProvider = rawProviderPodNotifierClient.(*podNotifier)

		return &providerWithPodNotifier{Provider:providerPlugin, PodNotifier: podNotifierProvider}, nil
	}

	return simpleProvider{Provider: providerPlugin}, nil
}

func loadConfig(providerConfig string, config *Configuration) error {
	if pluginName := os.Getenv("PLUGIN_NAME"); pluginName != "" {
		config.PluginName = pluginName
	} else if providerConfig != "" {
		// If this isn't set, then just opt for defaults
		file, err := os.Open(providerConfig)
		if err != nil {
			return errors.Wrapf(err, "Cannot open configuration file '%s'", providerConfig)
		}

		decoder := json.NewDecoder(file)
		err = decoder.Decode(config)
		if err != nil {
			return err
		}
	}
	if config.PluginName == "" {
		config.PluginName = "default"
	}
	if !pluginNameRegex.MatchString(config.PluginName) {
		return fmt.Errorf("Plugin name '%s' does not match regex %s", config.PluginName, pluginNameRegex.String())
	}

	return nil
}
