package runtime

import (
	"context"
	"net"
	"time"

	"github.com/wabenet/dodo-core/pkg/config"
	"google.golang.org/grpc"
	"k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const connectionTimout = 15 * time.Second

var _ Client = &client{}

type Client interface {
	cri.RuntimeService
	cri.ImageManagerService
}

type client struct {
	timeout       time.Duration
	runtimeClient runtimeapi.RuntimeServiceClient
	imageClient   runtimeapi.ImageServiceClient
}

func NewClient() (Client, error) {
	conn, err := grpc.Dial(
		config.GetCRIEndpoint(),
		grpc.WithInsecure(),
		grpc.WithTimeout(connectionTimout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		timeout:       connectionTimout,
		runtimeClient: runtimeapi.NewRuntimeServiceClient(conn),
		imageClient:   runtimeapi.NewImageServiceClient(conn),
	}, nil
}

// RuntimeVersioner

func (c *client) Version(ctx context.Context, apiVersion string) (*runtimeapi.VersionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.runtimeClient.Version(ctx, &runtimeapi.VersionRequest{Version: apiVersion})
}

// ContainerManager

func (c *client) CreateContainer(ctx context.Context, podSandBoxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.runtimeClient.CreateContainer(ctx, &runtimeapi.CreateContainerRequest{
		PodSandboxId:  podSandBoxID,
		Config:        config,
		SandboxConfig: sandboxConfig,
	})
	if err != nil {
		return "", err
	}

	return resp.ContainerId, nil
}

func (c *client) StartContainer(ctx context.Context, containerID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.runtimeClient.StartContainer(ctx, &runtimeapi.StartContainerRequest{ContainerId: containerID})

	return err
}

func (c *client) StopContainer(ctx context.Context, containerID string, timeout int64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	_, err := c.runtimeClient.StopContainer(ctx, &runtimeapi.StopContainerRequest{
		ContainerId: containerID,
		Timeout:     timeout,
	})

	return err
}

func (c *client) RemoveContainer(ctx context.Context, containerID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.runtimeClient.RemoveContainer(ctx, &runtimeapi.RemoveContainerRequest{ContainerId: containerID})

	return err
}

func (c *client) ListContainers(ctx context.Context, filter *runtimeapi.ContainerFilter) ([]*runtimeapi.Container, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.runtimeClient.ListContainers(ctx, &runtimeapi.ListContainersRequest{Filter: filter})
	if err != nil {
		return nil, err
	}

	return resp.Containers, nil
}

func (c *client) ContainerStatus(ctx context.Context, containerID string, verbose bool) (*runtimeapi.ContainerStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.runtimeClient.ContainerStatus(ctx, &runtimeapi.ContainerStatusRequest{ContainerId: containerID, Verbose: verbose})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) UpdateContainerResources(ctx context.Context, containerID string, resources *runtimeapi.ContainerResources) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.runtimeClient.UpdateContainerResources(ctx, &runtimeapi.UpdateContainerResourcesRequest{
		ContainerId: containerID,
		Linux:       resources.Linux,
		Windows:     resources.Windows,
	})

	return err
}

func (c *client) ExecSync(ctx context.Context, containerID string, cmd []string, timeout time.Duration) (stdout []byte, stderr []byte, err error) {
	var cancel context.CancelFunc

	if timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout+timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	defer cancel()

	resp, err := c.runtimeClient.ExecSync(ctx, &runtimeapi.ExecSyncRequest{
		ContainerId: containerID,
		Cmd:         cmd,
		Timeout:     int64(timeout.Seconds()),
	})
	if err != nil {
		return nil, nil, err
	}

	return resp.Stdout, resp.Stderr, err
}

func (c *client) Exec(ctx context.Context, req *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.runtimeClient.Exec(ctx, req)
}

func (c *client) Attach(ctx context.Context, req *runtimeapi.AttachRequest) (*runtimeapi.AttachResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.runtimeClient.Attach(ctx, req)
}

func (c *client) ReopenContainerLog(ctx context.Context, containerID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.runtimeClient.ReopenContainerLog(ctx, &runtimeapi.ReopenContainerLogRequest{ContainerId: containerID})

	return err
}

func (c *client) CheckpointContainer(ctx context.Context, req *runtimeapi.CheckpointContainerRequest) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.runtimeClient.CheckpointContainer(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) GetContainerEvents(ctx context.Context, containerEventsCh chan *runtimeapi.ContainerEventResponse, connectionEstablishedCallback func(runtimeapi.RuntimeService_GetContainerEventsClient)) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.runtimeClient.GetContainerEvents(ctx, &runtimeapi.GetEventsRequest{})
	if err != nil {
		return err
	}

	return nil
}

// PodSandboxManager

func (c *client) RunPodSandbox(ctx context.Context, config *runtimeapi.PodSandboxConfig, runtimeHandler string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout*2)
	defer cancel()

	resp, err := c.runtimeClient.RunPodSandbox(ctx, &runtimeapi.RunPodSandboxRequest{
		Config:         config,
		RuntimeHandler: runtimeHandler,
	})
	if err != nil {
		return "", err
	}

	return resp.PodSandboxId, nil
}

func (c *client) StopPodSandbox(ctx context.Context, podSandBoxID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.runtimeClient.StopPodSandbox(ctx, &runtimeapi.StopPodSandboxRequest{PodSandboxId: podSandBoxID})

	return err
}

func (c *client) RemovePodSandbox(ctx context.Context, podSandBoxID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.runtimeClient.RemovePodSandbox(ctx, &runtimeapi.RemovePodSandboxRequest{PodSandboxId: podSandBoxID})

	return err
}

func (c *client) PodSandboxStatus(ctx context.Context, podSandBoxID string, verbose bool) (*runtimeapi.PodSandboxStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.runtimeClient.PodSandboxStatus(ctx, &runtimeapi.PodSandboxStatusRequest{PodSandboxId: podSandBoxID, Verbose: verbose})
}

func (c *client) ListPodSandbox(ctx context.Context, filter *runtimeapi.PodSandboxFilter) ([]*runtimeapi.PodSandbox, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.runtimeClient.ListPodSandbox(ctx, &runtimeapi.ListPodSandboxRequest{Filter: filter})
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (c *client) PortForward(ctx context.Context, req *runtimeapi.PortForwardRequest) (*runtimeapi.PortForwardResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.runtimeClient.PortForward(ctx, req)
}

// ContainerStatsManager

func (c *client) ContainerStats(ctx context.Context, containerID string) (*runtimeapi.ContainerStats, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.runtimeClient.ContainerStats(ctx, &runtimeapi.ContainerStatsRequest{ContainerId: containerID})
	if err != nil {
		return nil, err
	}

	return resp.Stats, nil
}

func (c *client) ListContainerStats(ctx context.Context, filter *runtimeapi.ContainerStatsFilter) ([]*runtimeapi.ContainerStats, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resp, err := c.runtimeClient.ListContainerStats(ctx, &runtimeapi.ListContainerStatsRequest{Filter: filter})
	if err != nil {
		return nil, err
	}

	return resp.Stats, nil
}

func (c *client) PodSandboxStats(ctx context.Context, podSandboxID string) (*runtimeapi.PodSandboxStats, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.runtimeClient.PodSandboxStats(ctx, &runtimeapi.PodSandboxStatsRequest{PodSandboxId: podSandboxID})
	if err != nil {
		return nil, err
	}

	return resp.Stats, nil
}

func (c *client) ListPodSandboxStats(ctx context.Context, filter *runtimeapi.PodSandboxStatsFilter) ([]*runtimeapi.PodSandboxStats, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.runtimeClient.ListPodSandboxStats(ctx, &runtimeapi.ListPodSandboxStatsRequest{Filter: filter})
	if err != nil {
		return nil, err
	}

	return resp.Stats, nil
}

func (c *client) ListMetricDescriptors(ctx context.Context) ([]*runtimeapi.MetricDescriptor, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.runtimeClient.ListMetricDescriptors(ctx, &runtimeapi.ListMetricDescriptorsRequest{})
	if err != nil {
		return nil, err
	}

	return resp.Descriptors, nil
}

func (c *client) ListPodSandboxMetrics(ctx context.Context) ([]*runtimeapi.PodSandboxMetrics, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.runtimeClient.ListPodSandboxMetrics(ctx, &runtimeapi.ListPodSandboxMetricsRequest{})
	if err != nil {
		return nil, err
	}

	return resp.PodMetrics, nil
}

// RuntimeService

func (c *client) UpdateRuntimeConfig(ctx context.Context, runtimeConfig *runtimeapi.RuntimeConfig) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.runtimeClient.UpdateRuntimeConfig(ctx, &runtimeapi.UpdateRuntimeConfigRequest{
		RuntimeConfig: runtimeConfig,
	})

	return err
}

func (c *client) Status(ctx context.Context, verbose bool) (*runtimeapi.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.runtimeClient.Status(ctx, &runtimeapi.StatusRequest{Verbose: verbose})
}

func (c *client) RuntimeConfig(ctx context.Context) (*runtimeapi.RuntimeConfigResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.runtimeClient.RuntimeConfig(ctx, &runtimeapi.RuntimeConfigRequest{})
}

// ImageManagerService

func (c *client) ListImages(ctx context.Context, filter *runtimeapi.ImageFilter) ([]*runtimeapi.Image, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.imageClient.ListImages(ctx, &runtimeapi.ListImagesRequest{Filter: filter})
	if err != nil {
		return nil, err
	}

	return resp.Images, nil
}

func (c *client) ImageStatus(ctx context.Context, image *runtimeapi.ImageSpec, verbose bool) (*runtimeapi.ImageStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.imageClient.ImageStatus(ctx, &runtimeapi.ImageStatusRequest{Image: image, Verbose: verbose})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) PullImage(ctx context.Context, image *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resp, err := c.imageClient.PullImage(ctx, &runtimeapi.PullImageRequest{
		Image:         image,
		Auth:          auth,
		SandboxConfig: podSandboxConfig,
	})
	if err != nil {
		return "", err
	}

	return resp.ImageRef, nil
}

func (c *client) RemoveImage(ctx context.Context, image *runtimeapi.ImageSpec) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	_, err := c.imageClient.RemoveImage(ctx, &runtimeapi.RemoveImageRequest{Image: image})

	return err
}

func (c *client) ImageFsInfo(ctx context.Context) (*runtimeapi.ImageFsInfoResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return c.imageClient.ImageFsInfo(ctx, &runtimeapi.ImageFsInfoRequest{})
}
