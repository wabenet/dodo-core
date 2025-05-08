package plugin

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	api "github.com/wabenet/dodo-core/api/plugin/v1alpha1"
	"github.com/wabenet/dodo-core/pkg/config"
)

const (
	ProtocolVersion  = 1
	MagicCookieKey   = "DODO_PLUGIN"
	MagicCookieValue = "69318785-d741-4150-ac91-8f03fa703530"
)

type Plugin interface {
	PluginInfo() *api.PluginInfo
	Init() (Config, error)
	Cleanup()

	Type() Type
}

type Config map[string]string

type Type interface {
	String() string
	GRPCClient() (plugin.Plugin, error)
	GRPCServer(p Plugin) (plugin.Plugin, error)
}

type StreamConfig struct {
	Stdin          io.Reader
	Stdout         io.Writer
	Stderr         io.Writer
	TerminalHeight uint32
	TerminalWidth  uint32
}

type Manager struct {
	pluginTypes map[string]Type
	plugins     map[string]map[string]Plugin
}

func MkName(t Type, name string) *api.PluginName {
	pn := &api.PluginName{}

	pn.SetType(t.String())
	pn.SetName(name)

	return pn
}

func MkInfo(t Type, name string) *api.PluginInfo {
	pi := &api.PluginInfo{}

	pi.SetName(MkName(t, name))

	return pi
}

func Init() Manager {
	config.Configure()

	if os.Getenv(MagicCookieKey) == MagicCookieValue {
		log.SetDefault(log.New(config.GetPluginLoggerOptions()))
	} else {
		log.SetDefault(log.New(config.GetLoggerOptions()))
	}

	return Manager{
		pluginTypes: map[string]Type{},
		plugins:     map[string]map[string]Plugin{},
	}
}

func (m Manager) RegisterPluginTypes(ts ...Type) {
	for _, t := range ts {
		m.pluginTypes[t.String()] = t
	}
}

func (m Manager) IncludePlugins(ps ...Plugin) {
	for _, p := range ps {
		name := p.PluginInfo().GetName()

		if m.plugins[name.GetType()] == nil {
			m.plugins[name.GetType()] = map[string]Plugin{}
		}

		m.plugins[name.GetType()][name.GetName()] = p
	}
}

func (m Manager) ServePlugins(plugins ...Plugin) error {
	pluginMap := map[string]plugin.Plugin{}

	for _, p := range plugins {
		s, err := p.Type().GRPCServer(p)
		if err != nil {
			return fmt.Errorf("could not instantiate GRPC Server: %w", err)
		}

		pluginMap[p.PluginInfo().GetName().GetType()] = s
	}

	plugin.Serve(&plugin.ServeConfig{
		GRPCServer: plugin.DefaultGRPCServer,
		Plugins:    pluginMap,
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  ProtocolVersion,
			MagicCookieKey:   MagicCookieKey,
			MagicCookieValue: MagicCookieValue,
		},
	})

	return nil
}

func (m Manager) GetPlugins(pluginType string) map[string]Plugin {
	return m.plugins[pluginType]
}

func PathByName(name string) string {
	return filepath.Join(
		config.GetPluginDir(),
		fmt.Sprintf("dodo-%s_%s_%s", name, runtime.GOOS, runtime.GOARCH),
	)
}

func (m Manager) LoadPlugins() {
	m.findPlugins()

	ps := ResolveDependencies(m.plugins)

	for _, p := range ps {
		m.initPlugin(p)
	}
}

func (m Manager) UnloadPlugins() {
	for _, ps := range m.plugins {
		for _, p := range ps {
			p.Cleanup()
		}
	}

	plugin.CleanupClients()
}

func (m Manager) findPlugins() {
	matches, err := filepath.Glob(PathByName("*"))
	if err != nil {
		return
	}

	for _, path := range matches {
		logger := log.Default().With("path", path)

		if stat, err := os.Stat(path); err != nil || stat.Mode().Perm()&0o111 == 0 {
			continue
		}

		for n, t := range m.pluginTypes {
			logger.Debug("attempt loading plugin", "type", n)

			grpcClient, err := t.GRPCClient()
			if err != nil {
				logger.Debug("error loading plugin", "error", err)

				continue
			}

			p, err := loadGRPCPlugin(path, n, grpcClient)
			if err != nil {
				logger.Debug("could not load plugin over grpc", "error", err)

				continue
			}

			name := p.PluginInfo().GetName()

			if m.plugins[name.GetType()] == nil {
				m.plugins[name.GetType()] = map[string]Plugin{}
			}

			m.plugins[name.GetType()][name.GetName()] = p
		}
	}
}

func (m Manager) initPlugin(p Plugin) {
	info := p.PluginInfo()
	logger := log.L().With("name", info.GetName().GetName(), "type", info.GetName().GetType())
	logger = augmentLogger(logger, info.GetFields())

	if config, err := p.Init(); err != nil {
		logger.Warn("could not load plugin", "error", err)
		delete(m.plugins[info.GetName().GetType()], info.GetName().GetName())
	} else {
		augmentLogger(logger, config).Info("loaded plugin")
	}
}

func loadGRPCPlugin(path, pluginType string, grpcPlugin plugin.Plugin) (Plugin, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		Managed:          true,
		Plugins:          map[string]plugin.Plugin{pluginType: grpcPlugin},
		Cmd:              exec.Command(path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		Logger:           log.New(config.GetLoggerOptions()),
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  ProtocolVersion,
			MagicCookieKey:   MagicCookieKey,
			MagicCookieValue: MagicCookieValue,
		},
	})

	conn, err := client.Client()
	if err != nil {
		client.Kill()

		return nil, fmt.Errorf("error getting plugin client: %w", err)
	}

	raw, err := conn.Dispense(pluginType)
	if err != nil {
		client.Kill()

		return nil, fmt.Errorf("error dispensing plugin: %w", err)
	}

	if p, ok := raw.(Plugin); ok {
		return p, nil
	}

	client.Kill()

	invalid := &api.PluginName{}

	invalid.SetType(pluginType)
	invalid.SetName(path) // TODO: name?

	return nil, InvalidError{
		Plugin:  invalid,
		Message: "does not implement Plugin interface",
	}
}

func augmentLogger(logger log.Logger, fields map[string]string) log.Logger {
	fs := []interface{}{}

	for k, v := range fields {
		fs = append(fs, k, v)
	}

	return logger.With(fs...)
}
