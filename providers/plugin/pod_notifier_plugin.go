package plugin

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"github.com/virtual-kubelet/virtual-kubelet/providers"
	"github.com/virtual-kubelet/virtual-kubelet/providers/plugin/proto"
	"google.golang.org/grpc"
	"k8s.io/api/core/v1"
)

var (
	_ providers.PodNotifier = (*podNotifier)(nil)
	_ plugin.GRPCPlugin = (*podNotifierPlugin)(nil)
)

type podNotifierPlugin struct {
	plugin.NetRPCUnsupportedPlugin
}


func (*podNotifierPlugin) GRPCServer(*plugin.GRPCBroker, *grpc.Server) error {
	panic("This should not be invoked?")
}

func (*podNotifierPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, clientConn *grpc.ClientConn) (interface{}, error) {
	return &podNotifier{client: proto.NewPodNotifierProviderClient(clientConn), broker: broker}, nil
}

type podNotifier struct {
	client proto.PodNotifierProviderClient
	broker *plugin.GRPCBroker
}

func (p *podNotifier) NotifyPods(ctx context.Context, callback func(*v1.Pod)) {
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s := grpc.NewServer(opts...)
		proto.RegisterPodNotifierCallbackServer(s, &podNotifierCallbackServer{callback: callback})

		return s
	}

	brokerID := p.broker.NextId()
	go p.broker.AcceptAndServe(brokerID, serverFunc)
}


type podNotifierCallbackServer struct {
	callback func(*v1.Pod)
}

func (p *podNotifierCallbackServer) NotifyPods(ctx context.Context, notifyPodsRequest *proto.NotifyPodsRequest) (*proto.NotifyPodsResponse, error) {
	p.callback(notifyPodsRequest.GetPod())

	return &proto.NotifyPodsResponse{}, nil
}

