package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/wabenet/dodo-core/pkg/config"
)

const (
	ProtocolVersion  = 1
	MagicCookieKey   = "DODO_PLUGIN"
	MagicCookieValue = "69318785-d741-4150-ac91-8f03fa703530"
)

type Manager struct {
	pluginTypes map[string]Type
	plugins     map[string]map[string]Plugin
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

func (m Manager) IncludePlugins(plugins ...Plugin) {
	for _, plugin := range plugins {
		id := plugin.Metadata().ID

		if m.plugins[id.Type] == nil {
			m.plugins[id.Type] = map[string]Plugin{}
		}

		m.plugins[id.Type][id.Name] = plugin
	}
}

func (m Manager) ServePlugins(plugins ...Plugin) error {
	pluginMap := map[string]plugin.Plugin{}

	for _, plugin := range plugins {
		server, err := plugin.Type().GRPCServer(plugin)
		if err != nil {
			return fmt.Errorf("could not instantiate GRPC Server: %w", err)
		}

		pluginMap[plugin.Metadata().ID.Type] = server
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

		for name, pluginType := range m.pluginTypes {
			logger.Debug("attempt loading plugin", "type", pluginType)

			grpcClient, err := pluginType.GRPCClient()
			if err != nil {
				logger.Debug("error loading plugin", "error", err)

				continue
			}

			plugin, err := loadGRPCPlugin(path, name, grpcClient)
			if err != nil {
				logger.Debug("could not load plugin over grpc", "error", err)

				continue
			}

			id := plugin.Metadata().ID

			if m.plugins[id.Type] == nil {
				m.plugins[id.Type] = map[string]Plugin{}
			}

			m.plugins[id.Type][id.Name] = plugin
		}
	}
}

func (m Manager) initPlugin(p Plugin) {
	metadata := p.Metadata()
	logger := log.L().With("name", metadata.ID.Name, "type", metadata.ID.Type)
	logger = augmentLogger(logger, metadata.Labels)

	if config, err := p.Init(); err != nil {
		logger.Warn("could not load plugin", "error", err)
		delete(m.plugins[metadata.ID.Type], metadata.ID.Name)
	} else {
		augmentLogger(logger, config).Info("loaded plugin")
	}
}

func loadGRPCPlugin(path, pluginType string, grpcPlugin plugin.Plugin) (Plugin, error) { //nolint:ireturn
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

	return nil, InvalidError{
		PluginID: ID{Type: pluginType, Name: path}, // TODO: name?
		Message:  "does not implement Plugin interface",
	}
}

func augmentLogger(logger log.Logger, fields map[string]string) log.Logger { //nolint:ireturn
	all := []interface{}{}

	for k, v := range fields {
		all = append(all, k, v)
	}

	return logger.With(all...)
}
