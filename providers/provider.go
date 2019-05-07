package providers

import (
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
)

// Provider contains the methods required to implement a virtual-kubelet provider.
type Provider interface {
	// CreatePod takes a Kubernetes Pod and deploys it within the provider.
	CreatePod(ctx context.Context, pod *v1.Pod) error

	// UpdatePod takes a Kubernetes Pod and updates it within the provider.
	UpdatePod(ctx context.Context, pod *v1.Pod) error

	// DeletePod takes a Kubernetes Pod and deletes it from the provider.
	DeletePod(ctx context.Context, pod *v1.Pod) error

	// GetPod retrieves a pod by name from the provider (can be cached).
	GetPod(ctx context.Context, namespace, name string) (*v1.Pod, error)

	// GetContainerLogs retrieves the logs of a container by name from the provider.
	GetContainerLogs(ctx context.Context, namespace, podName, containerName string, tail int) (string, error)

	// RunInContainer executes a command in a container in the pod, copying data
	// between in/out/err and the container's stdin/stdout/stderr.
	RunInContainer(ctx context.Context, namespace, podName, containerName string, cmd []string, attach AttachIO) error

	// GetPodStatus retrieves the status of a pod by name from the provider.
	GetPodStatus(ctx context.Context, namespace, name string) (*v1.PodStatus, error)

	// GetPods retrieves a list of all pods running on the provider (can be cached).
	GetPods(context.Context) ([]*v1.Pod, error)

	// Capacity returns a resource list with the capacity constraints of the provider.
	Capacity(context.Context) v1.ResourceList

	// NodeConditions returns a list of conditions (Ready, OutOfDisk, etc), which is
	// polled periodically to update the node status within Kubernetes.
	NodeConditions(context.Context) []v1.NodeCondition

	// NodeAddresses returns a list of addresses for the node status
	// within Kubernetes.
	NodeAddresses(context.Context) []v1.NodeAddress

	// NodeDaemonEndpoints returns NodeDaemonEndpoints for the node status
	// within Kubernetes.
	NodeDaemonEndpoints(context.Context) *v1.NodeDaemonEndpoints

	// OperatingSystem returns the operating system the provider is for.
	OperatingSystem() string
}

// PodMetricsProvider is an optional interface that providers can implement to expose pod stats
type PodMetricsProvider interface {
	GetStatsSummary(context.Context) (*stats.Summary, error)
}

// PodNotifier notifies callers of pod changes.
// Providers should implement this interface to enable callers to be notified
// of pod status updates asyncronously.
type PodNotifier interface {
	// NotifyPods instructs the notifier to call the passed in function when
	// the pod status changes.
	//
	// NotifyPods should not block callers.
	NotifyPods(context.Context, func(*v1.Pod))
}

// SoftDeletes makes it so that containers are not immediately deleted from API server upon termination
// It is useful for providers that intend to have a terminal status update for a container.
// All Soft Delete providers must also implement the PodNotifier interface.
//
// You must implement this if you intend to ever delete pods after deletion
type SoftDeletes interface {
	PodNotifier
	// This method indicates the provider is okay with soft-deletes. Return true is the provider can handle them.
	SoftDeletes() bool
}

// AttachIO is used to pass in streams to attach to a container process
type AttachIO interface {
	Stdin() io.Reader
	Stdout() io.WriteCloser
	Stderr() io.WriteCloser
	TTY() bool
	Resize() <-chan TermSize
}

// TermSize is used to set the terminal size from attached clients.
type TermSize struct {
	Width  uint16
	Height uint16
}
