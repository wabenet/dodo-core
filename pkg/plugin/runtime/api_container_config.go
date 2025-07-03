package runtime

import (
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/runtime/v1alpha2"
)

type ContainerConfig struct {
	Name  string
	Image string

	Process  Process
	Terminal TerminalConfig

	Environment  EnvironmentConfig
	Ports        PortConfig
	Mounts       MountConfig
	Capabilities CapabilityConfig
}

func EmptyContainerConfig() ContainerConfig {
	return ContainerConfig{
		Process:      Process{},
		Terminal:     TerminalConfig{},
		Environment:  EnvironmentConfig{},
		Ports:        PortConfig{},
		Mounts:       MountConfig{},
		Capabilities: CapabilityConfig{},
	}
}

func MergeContainerConfig(first, second ContainerConfig) ContainerConfig {
	result := first

	if len(second.Name) > 0 {
		result.Name = second.Name
	}

	if len(second.Image) > 0 {
		result.Image = second.Image
	}

	result.Process = MergeProcess(first.Process, second.Process)
	result.Terminal = MergeTerminalConfig(first.Terminal, second.Terminal)
	result.Environment = MergeEnvironmentConfig(first.Environment, second.Environment)
	result.Ports = MergePortConfig(first.Ports, second.Ports)
	result.Mounts = MergeMountConfig(first.Mounts, second.Mounts)
	result.Capabilities = MergeCapabilityConfig(first.Capabilities, second.Capabilities)

	return result
}

func ContainerConfigFromProto(c *api.ContainerConfig) ContainerConfig {
	return ContainerConfig{
		Name:         c.GetName(),
		Image:        c.GetImage(),
		Process:      ProcessFromProto(c.GetProcess()),
		Terminal:     TerminalConfigFromProto(c.GetTerminal()),
		Environment:  EnvironmentConfigFromProto(c.GetEnvironment()),
		Ports:        PortConfigFromProto(c.GetPorts()),
		Mounts:       MountConfigFromProto(c.GetMounts()),
		Capabilities: c.GetCapabilities(),
	}
}

func (c ContainerConfig) ToProto() *api.ContainerConfig {
	out := &api.ContainerConfig{}

	out.SetName(c.Name)
	out.SetImage(c.Image)
	out.SetProcess(c.Process.ToProto())
	out.SetTerminal(c.Terminal.ToProto())
	out.SetEnvironment(c.Environment.ToProto())
	out.SetPorts(c.Ports.ToProto())
	out.SetMounts(c.Mounts.ToProto())
	out.SetCapabilities(c.Capabilities)

	return out
}

type Process struct {
	User       string
	WorkingDir string
	Entrypoint Entrypoint
	Command    Command
}

func MergeProcess(first, second Process) Process {
	result := first

	if len(second.User) > 0 {
		result.User = second.User
	}

	if len(second.WorkingDir) > 0 {
		result.WorkingDir = second.WorkingDir
	}

	if len(second.Entrypoint) > 0 {
		result.Entrypoint = second.Entrypoint
	}

	if len(second.Command) > 0 {
		result.Command = second.Command
	}

	return result
}

func ProcessFromProto(p *api.Process) Process {
	return Process{
		User:       p.GetUser(),
		WorkingDir: p.GetWorkingDir(),
		Entrypoint: p.GetEntrypoint(),
		Command:    p.GetCommand(),
	}
}

func (p Process) ToProto() *api.Process {
	out := &api.Process{}

	out.SetUser(p.User)
	out.SetWorkingDir(p.WorkingDir)
	out.SetEntrypoint(p.Entrypoint)
	out.SetCommand(p.Command)

	return out
}

type Entrypoint []string

type Command []string

type TerminalConfig struct {
	StdIO         bool
	TTY           bool
	ConsoleHeight int
	ConsoleWidth  int
}

func MergeTerminalConfig(first, second TerminalConfig) TerminalConfig {
	result := first

	if second.StdIO {
		result.StdIO = true
	}

	if second.TTY {
		result.TTY = true
	}

	if second.ConsoleHeight > 0 {
		result.ConsoleHeight = second.ConsoleHeight
	}

	if second.ConsoleWidth > 0 {
		result.ConsoleWidth = second.ConsoleWidth
	}

	return result
}

func TerminalConfigFromProto(t *api.TerminalConfig) TerminalConfig {
	return TerminalConfig{
		StdIO:         t.GetStdio(),
		TTY:           t.GetTty(),
		ConsoleHeight: int(t.GetConsoleHeight()),
		ConsoleWidth:  int(t.GetConsoleWidth()),
	}
}

func (t TerminalConfig) ToProto() *api.TerminalConfig {
	out := &api.TerminalConfig{}

	out.SetStdio(t.StdIO)
	out.SetTty(t.TTY)
	out.SetConsoleHeight(int64(t.ConsoleHeight))
	out.SetConsoleWidth(int64(t.ConsoleWidth))

	return out
}

type CapabilityConfig []string

func MergeCapabilityConfig(first, second CapabilityConfig) CapabilityConfig {
	return append(first, second...)
}
