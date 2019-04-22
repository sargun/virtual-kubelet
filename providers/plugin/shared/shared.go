package shared

import (
	"github.com/hashicorp/go-plugin"
)

const (
	magicCookieKey      = "VIRTUAL_KUBELET_PLUGIN"
	ProviderPluginName  = "provider"
	PodNotifierProviderName = "podNotifierProvider"
)

func HandshakeConfig(pluginName string) plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey: magicCookieKey,
		MagicCookieValue: pluginName,
		// We can switch this to using plugins with versions in v2
		ProtocolVersion: 1,
	}
}

