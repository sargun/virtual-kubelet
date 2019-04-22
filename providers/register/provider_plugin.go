// +build plugin_provider

package register

import (
	"github.com/virtual-kubelet/virtual-kubelet/providers"
	"github.com/virtual-kubelet/virtual-kubelet/providers/plugin"
)

func init() {
	register("plugin", initPlugin)
}

func initPlugin(cfg InitConfig) (providers.Provider, error) {
	return plugin.NewPluginProvider(
		cfg.ConfigPath,
		cfg.NodeName,
		cfg.OperatingSystem,
		cfg.InternalIP,
		cfg.DaemonPort,
	)
}