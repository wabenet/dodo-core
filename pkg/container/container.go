package container

import (
	"fmt"
	"os"
	"path/filepath"

	dockerapi "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/term"
	"github.com/oclaussen/dodo/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	DefaultAPIVersion = "1.39"
)

type ScriptError struct {
	Message  string
	ExitCode int
}

func (e *ScriptError) Error() string {
	return e.Message
}

type Container struct {
	name    string
	daemon  bool
	config  *types.Backdrop
	client  *client.Client
	context context.Context
	tmpPath string
}

func NewContainer(overrides *types.Backdrop, daemon bool) (*Container, error) {
	c := &Container{
		daemon: daemon,
		config: &types.Backdrop{
			Name:       overrides.Name,
			Entrypoint: &types.Entrypoint{},
		},
		context: context.Background(),
		tmpPath: fmt.Sprintf("/tmp/dodo-%s/", stringid.GenerateRandomID()[:20]),
	}

	for _, p := range plugin.GetPlugins(configuration.PluginType) {
		conf, err := p.(configuration.Configuration).UpdateConfiguration(c.config)
		if err != nil {
			log.Warn(err)
			continue
		}
		c.config = conf
	}
	c.config.Merge(overrides)

	log.WithFields(log.Fields{"backdrop": c.config}).Debug("assembled configuration")

	dockerClient, err := getDockerClient(c.config.Name)
	if err != nil {
		return nil, err
	}
	c.client = dockerClient

	c.name = c.config.ContainerName
	if c.daemon {
		c.name = c.config.Name
	} else if len(c.name) == 0 {
		c.name = fmt.Sprintf("%s-%s", c.config.Name, stringid.GenerateRandomID()[:8])
	}

	return c, nil
}

func (c *Container) Run() error {
	imageId, err := c.GetImage()
	if err != nil {
		return err
	}

	containerID, err := c.create(imageId)
	if err != nil {
		return err
	}

	if c.daemon {
		return c.client.ContainerStart(
			c.context,
			containerID,
			dockerapi.ContainerStartOptions{},
		)
	} else {
		return c.run(containerID, hasTTY())
	}
}

func (c *Container) Stop() error {
	if err := c.client.ContainerStop(c.context, c.name, nil); err != nil {
		return err
	}

	if err := c.client.ContainerRemove(c.context, c.name, dockerapi.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}

func hasTTY() bool {
	_, inTerm := term.GetFdInfo(os.Stdin)
	_, outTerm := term.GetFdInfo(os.Stdout)
	return inTerm && outTerm
}

func getDockerClient(name string) (*client.Client, error) {
	opts := &configuration.ClientOptions{
		Host: os.Getenv("DOCKER_HOST"),
	}
	if version := os.Getenv("DOCKER_API_VERSION"); len(version) > 0 {
		opts.Version = version
	}
	if certPath := os.Getenv("DOCKER_CERT_PATH"); len(certPath) > 0 {
		opts.CAFile = filepath.Join(certPath, "ca.pem")
		opts.CertFile = filepath.Join(certPath, "cert.pem")
		opts.KeyFile = filepath.Join(certPath, "key.pem")
	}

	for _, p := range plugin.GetPlugins(configuration.PluginType) {
		o, err := p.(configuration.Configuration).GetClientOptions(name)
		if err != nil {
			log.Warn(err)
			continue
		}
		if len(o.Host) > 0 { // FIXME: why only check host?
			opts = o
		}
	}

	mutators := []client.Opt{}
	if len(opts.Version) > 0 {
		mutators = append(mutators, client.WithVersion(opts.Version))
	} else {
		mutators = append(mutators, client.WithVersion(DefaultAPIVersion))
	}
	if len(opts.Host) > 0 {
		mutators = append(mutators, client.WithHost(opts.Host))
	}
	if len(opts.CAFile)+len(opts.CertFile)+len(opts.KeyFile) > 0 {
		mutators = append(mutators, client.WithTLSClientConfig(opts.CAFile, opts.CertFile, opts.KeyFile))
	}
	return client.NewClientWithOpts(mutators...)
}
