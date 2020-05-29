package types

func (target *Backdrop) Merge(source *Backdrop) {
	if len(source.Name) > 0 {
		target.Name = source.Name
	}

	target.Aliases = append(target.Aliases, source.Aliases...)

	if len(source.ImageId) > 0 {
		target.ImageId = source.ImageId
	}

	if source.Entrypoint != nil {
		if source.Entrypoint.Interactive {
			target.Entrypoint.Interactive = true
		}

		if len(source.Entrypoint.Interpreter) > 0 {
			target.Entrypoint.Interpreter = source.Entrypoint.Interpreter
		}

		if len(source.Entrypoint.Script) > 0 {
			target.Entrypoint.Script = source.Entrypoint.Script
		}

		if len(source.Entrypoint.Arguments) > 0 {
			target.Entrypoint.Arguments = source.Entrypoint.Arguments
		}
	}

	if len(source.ContainerName) > 0 {
		target.ContainerName = source.ContainerName
	}

	target.Environment = append(target.Environment, source.Environment...)

	if len(source.User) > 0 {
		target.User = source.User
	}

	target.Volumes = append(target.Volumes, source.Volumes...)
	target.Devices = append(target.Devices, source.Devices...)
	target.Ports = append(target.Ports, source.Ports...)

	if len(source.WorkingDir) > 0 {
		target.WorkingDir = source.WorkingDir
	}
}
