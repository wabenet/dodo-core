package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/oclaussen/dodo/pkg/types"
	"github.com/spf13/cobra"
)

type options struct {
	interactive bool
	user        string
	workdir     string
	volumes     []string
	environment []string
	publish     []string
}

func (opts *options) createFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.SetInterspersed(false)

	flags.BoolVarP(
		&opts.interactive, "interactive", "i", false,
		"run an interactive session")
	flags.StringVarP(
		&opts.user, "user", "u", "",
		"username or UID (format: <name|uid>[:<group|gid>])")
	flags.StringVarP(
		&opts.workdir, "workdir", "w", "",
		"working directory inside the container")
	flags.StringArrayVarP(
		&opts.volumes, "volume", "v", []string{},
		"bind mount a volume")
	flags.StringArrayVarP(
		&opts.environment, "env", "e", []string{},
		"set environment variables")
	flags.StringArrayVarP(
		&opts.publish, "publish", "p", []string{},
		"publish a container's port(s) to the host")
}

func (opts *options) createConfig(name string, command []string) (*types.Backdrop, error) {
	config := &types.Backdrop{
		Name: name,
		Entrypoint: &types.Entrypoint{
			Interactive: opts.interactive,
			Arguments:   command,
		},
		User:       opts.user,
		WorkingDir: opts.workdir,
	}

	for _, volume := range opts.volumes {
		var v *types.Volume
		switch values := strings.SplitN(volume, ":", 3); len(values) {
		case 1:
			v = &types.Volume{Source: values[0]}
		case 2:
			v = &types.Volume{
				Source: values[0],
				Target: values[1],
			}
		case 3:
			v = &types.Volume{
				Source:   values[0],
				Target:   values[1],
				Readonly: values[2] == "ro",
			}
		default:
			return nil, fmt.Errorf("invalid volume definition: %s", volume)
		}
		config.Volumes = append(config.Volumes, v)
	}

	for _, env := range opts.environment {
		var e *types.Environment
		switch values := strings.SplitN(env, "=", 2); len(values) {
		case 1:
			e = &types.Environment{Key: values[0], Value: os.Getenv(values[0])}
		case 2:
			e = &types.Environment{Key: values[0], Value: values[1]}
		default:
			return nil, fmt.Errorf("invalid environment definition: %s", env)
		}
		config.Environment = append(config.Environment, e)
	}

	for _, port := range opts.publish {
		var p *types.Port
		switch values := strings.SplitN(port, ":", 3); len(values) {
		case 1:
			p.Target = values[0]
		case 2:
			p.Published = values[0]
			p.Target = values[1]
		case 3:
			p.HostIp = values[0]
			p.Published = values[1]
			p.Target = values[2]
		default:
			return nil, fmt.Errorf("invalid publish definition: %s", port)
		}

		switch values := strings.SplitN(p.Target, "/", 2); len(values) {
		case 1:
			p.Target = values[0]
		case 2:
			p.Target = values[0]
			p.Protocol = values[1]
		default:
			return nil, fmt.Errorf("invalid publish definition: %s", port)
		}

		config.Ports = append(config.Ports, p)
	}

	return config, nil
}
