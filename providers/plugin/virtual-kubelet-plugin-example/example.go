package main

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"github.com/virtual-kubelet/virtual-kubelet/providers/mock"
	"github.com/virtual-kubelet/virtual-kubelet/providers/plugin/proto"
	"github.com/virtual-kubelet/virtual-kubelet/providers/plugin/shared"
	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
)

const (
	pluginName = "example"
)

var (
	_ plugin.GRPCPlugin = (*exampleProviderPlugin)(nil)
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig(pluginName),
		Plugins: map[string]plugin.Plugin{
			shared.ProviderPluginName: &exampleProviderPlugin{},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,

	})
}

type exampleProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin
}


func (p *exampleProviderPlugin) GRPCClient(context.Context, *plugin.GRPCBroker, *grpc.ClientConn) (interface{}, error) {
	panic("Plugin does not implement client")
}

func (p *exampleProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterProviderServer(s, &exampleProviderPluginServer{})
	return nil
}

type exampleProviderPluginServer struct {
	mockProvider *mock.MockProvider
}

func (e *exampleProviderPluginServer) Register(ctx context.Context, req *proto.ProviderRegisterRequest) (*proto.ProviderRegisterResponse, error) {
	var err error
	cfg := req.GetInitConfig()
	e.mockProvider, err = mock.NewMockProvider(cfg.GetConfigPath(), cfg.GetNodeName(), cfg.GetOperatingSystem(), cfg.GetInternalIP(), cfg.GetDaemonPort())

	return &proto.ProviderRegisterResponse{}, err
}

func (e *exampleProviderPluginServer) CreatePod(ctx context.Context, cpr *proto.CreatePodRequest) (*proto.CreatePodResponse, error) {
	return &proto.CreatePodResponse{}, e.mockProvider.CreatePod(ctx, cpr.GetPod())
}

func (e *exampleProviderPluginServer) UpdatePod(ctx context.Context, upr *proto.UpdatePodRequest) (*proto.UpdatePodResponse, error) {
	return &proto.UpdatePodResponse{}, e.mockProvider.CreatePod(ctx, upr.GetPod())
}

func (e *exampleProviderPluginServer) DeletePod(ctx context.Context, dpr *proto.DeletePodRequest) (*proto.DeletePodResponse, error) {
	return &proto.DeletePodResponse{}, e.mockProvider.DeletePod(ctx, dpr.GetPod())
}

func (e *exampleProviderPluginServer) GetPod(ctx context.Context, gpr *proto.GetPodRequest) (*proto.GetPodResponse, error) {
	pod, err := e.mockProvider.GetPod(ctx, gpr.GetNamespace(), gpr.GetName())
	return &proto.GetPodResponse{
		Pod:pod,
	}, err
}

func (e *exampleProviderPluginServer) GetContainerLogs(ctx context.Context, gclr *proto.GetContainerLogsRequest) (*proto.GetContainerLogsResponse, error) {
	containerLogs, err := e.mockProvider.GetContainerLogs(ctx, gclr.GetNamespace(), gclr.GetPodName(), gclr.GetContainerName(), int(gclr.GetTail()))
	return &proto.GetContainerLogsResponse{
		Logs:containerLogs,
	}, err
}

func (e *exampleProviderPluginServer) GetPodStatus(ctx context.Context, gpsr *proto.GetPodStatusRequest) (*proto.GetPodStatusResponse, error) {
	status, err := e.mockProvider.GetPodStatus(ctx, gpsr.GetNamespace(), gpsr.GetName())
	return &proto.GetPodStatusResponse{
		Status: status,
	}, err}

func (e *exampleProviderPluginServer) GetPods(ctx context.Context, gpr *proto.GetPodsRequest) (*proto.GetPodsResponse, error) {
	pods, err := e.mockProvider.GetPods(ctx)
	return &proto.GetPodsResponse{
		Pods:pods,
	}, err
	}

func (e *exampleProviderPluginServer) Capacity(ctx context.Context, cr *proto.CapacityRequest) (*proto.CapacityResponse, error) {
	capacity := e.mockProvider.Capacity(ctx)
	resourceListResponse := make(map[string]string, len(capacity))
	for resourceName, resource := range capacity {
		resourceListResponse[string(resourceName)] = resource.String()
	}
	return &proto.CapacityResponse{
		ResourceList: resourceListResponse,
	}, nil
}

func (e *exampleProviderPluginServer) NodeConditions(ctx context.Context, ncr *proto.NodeConditionsRequest) (*proto.NodeConditionsResponse, error) {
	nodeConditions := e.mockProvider.NodeConditions(ctx)
	nodeConditionsResponse := make([]*v1.NodeCondition, len(nodeConditions))
	for idx, nodeCondition := range nodeConditions {
		nodeConditionsResponse[idx] = &nodeCondition
	}
	return &proto.NodeConditionsResponse{
		NodeConditions: nodeConditionsResponse,
	}, nil
}

func (e *exampleProviderPluginServer) NodeAddresses(ctx context.Context, ner *proto.NodeAddressesRequest) (*proto.NodeAddressesResponse, error) {
	nodeAddresses := e.mockProvider.NodeAddresses(ctx)
	nodeAddressesResponse := make([]*v1.NodeAddress, len(nodeAddresses))
	for idx, nodeAddress := range nodeAddresses {
		nodeAddressesResponse[idx] = &nodeAddress
	}

	return &proto.NodeAddressesResponse{
		NodeAddresses: nodeAddressesResponse,
	}, nil
}

func (e *exampleProviderPluginServer) NodeDaemonEndspoints(ctx context.Context, nder *proto.NodeDaemonEndpointsRequest) (*proto.NodeDaemonEndpointsResponse, error) {
	return &proto.NodeDaemonEndpointsResponse{
		NodeDaemonEndpoints: e.mockProvider.NodeDaemonEndpoints(ctx),
	}, nil
}

func (e *exampleProviderPluginServer) OperatingSystem(ctx context.Context, osr *proto.OperatingSystemRequest) (*proto.OperatingSystemResponse, error) {
	os := e.mockProvider.OperatingSystem()
	return &proto.OperatingSystemResponse{OperatingSystem: os}, nil
}
