package types

import (
	"path"

	core "github.com/wabenet/dodo-core/api/core/v1alpha5"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func RuntimeConfig(b *core.Backdrop, tmpPath string, tty bool, stdio bool) (*runtimeapi.PodSandboxConfig, *runtimeapi.ContainerConfig) {
	sandbox := &runtimeapi.PodSandboxConfig{
		Metadata: &runtimeapi.PodSandboxMetadata{Name: b.ContainerName},
		Hostname: b.ContainerName,
	}

	sandbox.PortMappings = []*runtimeapi.PortMapping{}
	for _, port := range b.Ports {
		mapping := &runtimeapi.PortMapping{
			HostPort:      port.Published,
			ContainerPort: port.Target,
			HostIp:        port.HostIp,
		}



		sandbox.PortMappings = append(sandbox.PortMappings, mapping)
	}

	container := &runtimeapi.ContainerConfig{
		Metadata:   &runtimeapi.ContainerMetadata{Name: b.ContainerName},
		Image:      &runtimeapi.ImageSpec{Image: b.ImageId},
		WorkingDir: b.WorkingDir,
		Stdin:      stdio,
		StdinOnce:  stdio,
		Tty:        tty && stdio,
		Linux: &runtimeapi.LinuxContainerConfig{
			SecurityContext: &runtimeapi.LinuxContainerSecurityContext{
				Capabilities: &runtimeapi.Capability{
					AddCapabilities: b.Capabilities,
				},
				RunAsUser:          &runtimeapi.Int64Value{Value: 0}, // TODO
				RunAsUsername:      "",                               // TODO
				RunAsGroup:         &runtimeapi.Int64Value{Value: 0}, // TODO
				ReadonlyRootfs:     false,                            // TODO
				SupplementalGroups: []int64{},                        // TODO
			},
		},
	}

	if b.Entrypoint.Interpreter == nil {
		container.Command = []string{"/bin/sh"}
	} else {
		container.Command = b.Entrypoint.Interpreter
	}

	container.Args = b.Entrypoint.Arguments

	if b.Entrypoint.Interactive {
		container.Args = nil
	} else if len(b.Entrypoint.Script) > 0 {
		container.Command = append(container.Command, path.Join(tmpPath, "entrypoint"))
	}

	container.Envs = []*runtimeapi.KeyValue{}
	for _, kv := range b.Environment {
		container.Envs = append(container.Envs, &runtimeapi.KeyValue{
			Key:   kv.Key,
			Value: kv.Value,
		})
	}

	container.Mounts = []*runtimeapi.Mount{}
	for _, v := range b.Volumes {
		container.Mounts = append(container.Mounts, &runtimeapi.Mount{
			HostPath:      v.Source,
			ContainerPath: v.Target,
			Readonly:      v.Readonly,
		})
	}

	container.Devices = []*runtimeapi.Device{}
	for _, device := range b.Devices {
		container.Devices = append(container.Devices, &runtimeapi.Device{
			HostPath:      device.Source,
			ContainerPath: device.Target,
			Permissions:   device.Permissions,
		})
	}

	return sandbox, container
}
