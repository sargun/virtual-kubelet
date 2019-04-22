package plugin

import (
	"context"
	"errors"
	"github.com/hashicorp/go-plugin"
	"github.com/virtual-kubelet/virtual-kubelet/providers"
	"github.com/virtual-kubelet/virtual-kubelet/providers/plugin/proto"
	"google.golang.org/grpc"
	"io"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/remotecommand"
	"time"
)

var (
	errExecInContainerNotImplemented = errors.New("Plugin provider does not implement ExecInContainer")
	// Make sure ProviderPlugin implements provider interface
	_ providers.Provider = (*ProviderPlugin)(nil)
	_ plugin.GRPCPlugin = (*providerPlugin)(nil)
)

type providerPlugin struct {
	plugin.NetRPCUnsupportedPlugin

}

func (*providerPlugin) GRPCServer(*plugin.GRPCBroker, *grpc.Server) error {
	panic("This should not be invoked?")
}

func (*providerPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, clientConn *grpc.ClientConn) (interface{}, error) {
	return &ProviderPlugin{client: proto.NewProviderClient(clientConn)}, nil
}


type ProviderPlugin struct {
	client proto.ProviderClient
}

func (p *ProviderPlugin) CreatePod(ctx context.Context, pod *v1.Pod) error {
	_, e := p.client.CreatePod(ctx, &proto.CreatePodRequest{
		Pod: pod,
	})
	return e
}

func (p *ProviderPlugin) UpdatePod(ctx context.Context, pod *v1.Pod) error {
	_, e := p.client.UpdatePod(ctx, &proto.UpdatePodRequest{
		Pod: pod,
	})
	return e
}

func (p *ProviderPlugin) DeletePod(ctx context.Context, pod *v1.Pod) error {
	_, e := p.client.DeletePod(ctx, &proto.DeletePodRequest{
		Pod: pod,
	})
	return e
}

func (p *ProviderPlugin) GetPod(ctx context.Context, namespace, name string) (*v1.Pod, error) {
	getPodResponse, e := p.client.GetPod(ctx, &proto.GetPodRequest{
		Namespace: namespace,
		Name: name,
	})

	return getPodResponse.GetPod(), e
}

func (p *ProviderPlugin) GetContainerLogs(ctx context.Context, namespace, podName, containerName string, tail int) (string, error) {
	getContainerLogsResponse, e := p.client.GetContainerLogs(ctx, &proto.GetContainerLogsRequest{
		Namespace: namespace,
		PodName: podName,
		ContainerName: containerName,
	})

	return getContainerLogsResponse.GetLogs(), e
}

func (ProviderPlugin) ExecInContainer(name string, uid types.UID, container string, cmd []string, in io.Reader, out, err io.WriteCloser, tty bool, resize <-chan remotecommand.TerminalSize, timeout time.Duration) error {
	return errExecInContainerNotImplemented
}

func (p *ProviderPlugin) GetPodStatus(ctx context.Context, namespace, name string) (*v1.PodStatus, error) {
	getPodStatusResponse, e := p.client.GetPodStatus(ctx, &proto.GetPodStatusRequest{
		Namespace: namespace,
		Name: name,
	})

	return getPodStatusResponse.GetStatus(), e
}

func (p *ProviderPlugin) GetPods(ctx context.Context) ([]*v1.Pod, error) {
	getPodsResponse, e := p.client.GetPods(ctx, &proto.GetPodsRequest{})

	return getPodsResponse.GetPods(), e
}

func (p *ProviderPlugin) Capacity(ctx context.Context) v1.ResourceList {
	capacityResponse, e := p.client.Capacity(ctx, &proto.CapacityRequest{})

	// Is there a smarter thing to do here?
	if e != nil {
		panic( e)
	}

	resourceList := make(map[v1.ResourceName]resource.Quantity, len(capacityResponse.GetResourceList()))
	for resourceName, canonicalResourceValue := range capacityResponse.GetResourceList() {
		resourceList[v1.ResourceName(resourceName)] = resource.MustParse(canonicalResourceValue)
	}
	return resourceList
}

func (p *ProviderPlugin) NodeConditions(ctx context.Context) []v1.NodeCondition {
	nodeConditionsResponse, e := p.client.NodeConditions(ctx, &proto.NodeConditionsRequest{})

	// Is there a smarter thing to do here?
	if e != nil {
		panic( e)
	}

	nodeConditions := make([]v1.NodeCondition, len(nodeConditionsResponse.GetNodeConditions()))
	for idx, nodeCondition := range nodeConditionsResponse.GetNodeConditions() {
		nodeConditions[idx] = *nodeCondition
	}
	return nodeConditions
}

func (p *ProviderPlugin) NodeAddresses(ctx context.Context) []v1.NodeAddress {
	nodeAddressesResponse, e := p.client.NodeAddresses(ctx, &proto.NodeAddressesRequest{})

	// Is there a smarter thing to do here?
	if e != nil {
		panic( e)
	}

	nodeAddresses := make([]v1.NodeAddress, len(nodeAddressesResponse.GetNodeAddresses()))
	for idx, nodeCondition := range nodeAddressesResponse.GetNodeAddresses() {
		nodeAddresses[idx] = *nodeCondition
	}
	return nodeAddresses
}

func (p *ProviderPlugin) NodeDaemonEndpoints(ctx context.Context) *v1.NodeDaemonEndpoints {
	nodeDaemonEndspointResponse, e := p.client.NodeDaemonEndspoints(ctx, &proto.NodeDaemonEndpointsRequest{})

	// Is there a smarter thing to do here?
	if e != nil {
		panic(e)
	}
	return nodeDaemonEndspointResponse.GetNodeDaemonEndpoints()
}

func (p *ProviderPlugin) OperatingSystem() string {
	operatingSystemResponse, e := p.client.OperatingSystem(context.TODO(), &proto.OperatingSystemRequest{})

	// Is there a smarter thing to do here?
	if e != nil {
		panic(e)
	}

	return operatingSystemResponse.GetOperatingSystem()
}

