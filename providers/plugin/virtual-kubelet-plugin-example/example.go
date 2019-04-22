package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-plugin"
	"github.com/virtual-kubelet/virtual-kubelet/providers/mock"
	"github.com/virtual-kubelet/virtual-kubelet/providers/plugin/proto"
	"github.com/virtual-kubelet/virtual-kubelet/providers/plugin/shared"
	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
	"time"
)

const (
	pluginName = "example"
)

var (
	_ plugin.GRPCPlugin = (*exampleProviderPlugin)(nil)
)

func main() {
	p := &provider{}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig(pluginName),
		VersionedPlugins: map[int]plugin.PluginSet{
			shared.VersionWithFeatures(1): map[string]plugin.Plugin{
				shared.ProviderPluginName: &exampleProviderPlugin{provider: p},
			},
			shared.VersionWithFeatures(1, shared.PodNotifier): map[string]plugin.Plugin{
				shared.ProviderPluginName: &exampleProviderPlugin{provider: p},
				shared.PodNotifierPluginName: &examplePodNotifierProviderPlugin{provider: p},
			},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,

	})
}

type exampleProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	provider *provider
}

func (p *exampleProviderPlugin) GRPCClient(context.Context, *plugin.GRPCBroker, *grpc.ClientConn) (interface{}, error) {
	panic("Plugin does not implement client")
}

func (p *exampleProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterProviderServer(s, p.provider)
	return nil
}

type examplePodNotifierProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	provider *provider
}

func (p *examplePodNotifierProviderPlugin) GRPCClient(context.Context, *plugin.GRPCBroker, *grpc.ClientConn) (interface{}, error) {
	panic("Plugin does not implement client")
}

func (p *examplePodNotifierProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPodNotifierProviderServer(s, &podNotifyProvider{provider: p.provider, broker: broker})
	return nil
}

type provider struct {
	mockProvider *mock.MockProvider
}

func (e *provider) Register(ctx context.Context, req *proto.ProviderRegisterRequest) (*proto.ProviderRegisterResponse, error) {

	var err error
	cfg := req.GetInitConfig()
	e.mockProvider, err = mock.NewMockProvider(cfg.GetConfigPath(), cfg.GetNodeName(), cfg.GetOperatingSystem(), cfg.GetInternalIP(), cfg.GetDaemonPort())

	return &proto.ProviderRegisterResponse{}, err
}

func (e *provider) CreatePod(ctx context.Context, cpr *proto.CreatePodRequest) (*proto.CreatePodResponse, error) {
	fmt.Println("CreatePod")

	return &proto.CreatePodResponse{}, e.mockProvider.CreatePod(ctx, cpr.GetPod())
}

func (e *provider) UpdatePod(ctx context.Context, upr *proto.UpdatePodRequest) (*proto.UpdatePodResponse, error) {
	fmt.Println("UpdatePod")

	return &proto.UpdatePodResponse{}, e.mockProvider.CreatePod(ctx, upr.GetPod())
}

func (e *provider) DeletePod(ctx context.Context, dpr *proto.DeletePodRequest) (*proto.DeletePodResponse, error) {
	fmt.Println("DeletePod")

	return &proto.DeletePodResponse{}, e.mockProvider.DeletePod(ctx, dpr.GetPod())
}

func (e *provider) GetPod(ctx context.Context, gpr *proto.GetPodRequest) (*proto.GetPodResponse, error) {
	fmt.Println("GetPod")
	pod, err := e.mockProvider.GetPod(ctx, gpr.GetNamespace(), gpr.GetName())
	return &proto.GetPodResponse{
		Pod:pod,
	}, err
}

func (e *provider) GetContainerLogs(ctx context.Context, gclr *proto.GetContainerLogsRequest) (*proto.GetContainerLogsResponse, error) {
	containerLogs, err := e.mockProvider.GetContainerLogs(ctx, gclr.GetNamespace(), gclr.GetPodName(), gclr.GetContainerName(), int(gclr.GetTail()))
	return &proto.GetContainerLogsResponse{
		Logs:containerLogs,
	}, err
}

func (e *provider) GetPodStatus(ctx context.Context, gpsr *proto.GetPodStatusRequest) (*proto.GetPodStatusResponse, error) {
	fmt.Println("Get Pod Status")
	status, err := e.mockProvider.GetPodStatus(ctx, gpsr.GetNamespace(), gpsr.GetName())
	return &proto.GetPodStatusResponse{
		Status: status,
	}, err}

func (e *provider) GetPods(ctx context.Context, gpr *proto.GetPodsRequest) (*proto.GetPodsResponse, error) {
	fmt.Println("Get Pods")
	pods, err := e.mockProvider.GetPods(ctx)
	return &proto.GetPodsResponse{
		Pods:pods,
	}, err
	}

func (e *provider) Capacity(ctx context.Context, cr *proto.CapacityRequest) (*proto.CapacityResponse, error) {
	capacity := e.mockProvider.Capacity(ctx)
	resourceListResponse := make(map[string]string, len(capacity))
	for resourceName, resource := range capacity {
		resourceListResponse[string(resourceName)] = resource.String()
	}
	return &proto.CapacityResponse{
		ResourceList: resourceListResponse,
	}, nil
}

func (e *provider) NodeConditions(ctx context.Context, ncr *proto.NodeConditionsRequest) (*proto.NodeConditionsResponse, error) {
	nodeConditions := e.mockProvider.NodeConditions(ctx)
	nodeConditionsResponse := make([]*v1.NodeCondition, len(nodeConditions))
	for idx := range nodeConditions {
		nodeConditionsResponse[idx] = &nodeConditions[idx]
	}
	return &proto.NodeConditionsResponse{
		NodeConditions: nodeConditionsResponse,
	}, nil
}

func (e *provider) NodeAddresses(ctx context.Context, ner *proto.NodeAddressesRequest) (*proto.NodeAddressesResponse, error) {
	nodeAddresses := e.mockProvider.NodeAddresses(ctx)
	nodeAddressesResponse := make([]*v1.NodeAddress, len(nodeAddresses))
	for idx := range nodeAddresses {
		nodeAddressesResponse[idx] = &nodeAddresses[idx]
	}

	return &proto.NodeAddressesResponse{
		NodeAddresses: nodeAddressesResponse,
	}, nil
}

func (e *provider) NodeDaemonEndspoints(ctx context.Context, nder *proto.NodeDaemonEndpointsRequest) (*proto.NodeDaemonEndpointsResponse, error) {
	return &proto.NodeDaemonEndpointsResponse{
		NodeDaemonEndpoints: e.mockProvider.NodeDaemonEndpoints(ctx),
	}, nil
}

func (e *provider) OperatingSystem(ctx context.Context, osr *proto.OperatingSystemRequest) (*proto.OperatingSystemResponse, error) {
	os := e.mockProvider.OperatingSystem()
	return &proto.OperatingSystemResponse{OperatingSystem: os}, nil
}

type podNotifyProvider struct {
	provider *provider
	broker *plugin.GRPCBroker
}

func (p *podNotifyProvider) Register(ctx context.Context, req *proto.PodNotifierProviderRegisterRequest) (*proto.PodNotifierProviderRegisterResponse, error) {
	fmt.Println("Got pod notification registration")
	conn, err := p.broker.Dial(req.GetPodNotifierProviderBrokerId())
	if err != nil {
		return nil, err
	}

	client := proto.NewPodNotifierCallbackClient(conn)

	go p.notify(client)

	return &proto.PodNotifierProviderRegisterResponse{}, nil
}

func (p *podNotifyProvider) notify(client proto.PodNotifierCallbackClient) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for range ticker.C {
		fmt.Println("Doing the notification dance?")
		pods, err := p.provider.mockProvider.GetPods(ctx)
		if err != nil {
			panic(err)
		}
		for _, pod := range pods {
			client.NotifyPods(ctx, &proto.NotifyPodsRequest{Pod: pod}, grpc.CallContentSubtype(shared.JSONCodecName))
		}
	}
}