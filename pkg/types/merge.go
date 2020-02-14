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
		target.Entrypoint.Merge(source.Entrypoint)
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

func (target *Entrypoint) Merge(source *Entrypoint) {
	if source.Interactive {
		target.Interactive = true
	}
	if len(source.Interpreter) > 0 {
		target.Interpreter = source.Interpreter
	}
	if len(source.Script) > 0 {
		target.Script = source.Script
	}
	if len(source.Arguments) > 0 {
		target.Arguments = source.Arguments
	}
}
