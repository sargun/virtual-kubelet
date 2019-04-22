//go:generate protoc -I ./proto/ -I ../../vendor/ proto/plugin.proto --go_out=plugins=grpc:proto/

package plugin
