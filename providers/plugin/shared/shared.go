package shared

import (
	"encoding/json"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc/encoding"
)

var (
	_ encoding.Codec = (*grpcJSONCodec)(nil)
)

const (
	magicCookieKey        = "VIRTUAL_KUBELET_PLUGIN"
	ProviderPluginName    = "provider"
	PodNotifierPluginName = "podNotifierProvider"
	MetricsPluginName      =      "metricsProvider"
	JSONCodecName         = "json"
)

func HandshakeConfig(pluginName string) plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey: magicCookieKey,
		MagicCookieValue: pluginName,
		// We can switch this to using plugins with versions in v2
	}
}

type Feature uint16
type Subversion uint16

const (
	Metrics Feature = 1 << iota
	PodNotifier
)

func HasFeature(version int, feature Feature) bool {
	return uint16(0xffff & version) & uint16(feature) > 0
}

func VersionWithFeatures(version Subversion, features... Feature) int {
	// Go ints are at least 32 bits
	totalVersion := uint32(version) << 16

	for _, feature := range features {
		totalVersion = totalVersion | uint32(feature)
	}

	return int(totalVersion)
}

func init() {
	encoding.RegisterCodec(grpcJSONCodec{})
}

type grpcJSONCodec struct {

}

func (grpcJSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (grpcJSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (grpcJSONCodec) Name() string {
	return JSONCodecName
}

