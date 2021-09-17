package plugin

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/config"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

const (
	ProtocolVersion  = 1
	MagicCookieKey   = "DODO_PLUGIN"
	MagicCookieValue = "69318785-d741-4150-ac91-8f03fa703530"

	FailedPlugin = "error"

	ErrPluginInvalid        PluginError = "invalid plugin"
	ErrPluginNotImplemented PluginError = "not implemented"
	ErrPluginNotFound       PluginError = "plugin not found"
	ErrNoValidPluginFound   PluginError = "no valid plugin found"
	ErrCircularDependency   PluginError = "circular plugin dependency"
)

var (
	pluginTypes = map[string]Type{}
	plugins     = map[string]map[string]Plugin{}
)

type PluginError string

func (e PluginError) Error() string {
	return string(e)
}

type Plugin interface {
	PluginInfo() *api.PluginInfo
	Init() (PluginConfig, error)

	Type() Type
}

type PluginConfig map[string]string

type Type interface {
	String() string
	GRPCClient() (plugin.Plugin, error)
	GRPCServer(Plugin) (plugin.Plugin, error)
}

type StreamConfig struct {
	Stdin          io.Reader
	Stdout         io.Writer
	Stderr         io.Writer
	TerminalHeight uint32
	TerminalWidth  uint32
}

func RegisterPluginTypes(ts ...Type) {
	for _, t := range ts {
		pluginTypes[t.String()] = t
	}
}

func IncludePlugins(ps ...Plugin) {
	for _, p := range ps {
		name := p.PluginInfo().Name

		if plugins[name.Type] == nil {
			plugins[name.Type] = map[string]Plugin{}
		}

		plugins[name.Type][name.Name] = p
	}
}

func ServePlugins(plugins ...Plugin) error {
	log.SetDefault(log.New(config.GetPluginLoggerOptions()))

	pluginMap := map[string]plugin.Plugin{}

	for _, p := range plugins {
		s, err := p.Type().GRPCServer(p)
		if err != nil {
			return fmt.Errorf("could not instantiate GRPC Server: %w", err)
		}

		pluginMap[p.PluginInfo().Name.Type] = s
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

func GetPlugins(pluginType string) map[string]Plugin {
	return plugins[pluginType]
}

func PathByName(name string) string {
	return filepath.Join(
		config.GetPluginDir(),
		fmt.Sprintf("dodo-%s_%s_%s", name, runtime.GOOS, runtime.GOARCH),
	)
}

func LoadPlugins() {
	findPlugins()

	ps, err := ResolveDependencies(plugins)
	if err != nil {
		log.L().Error("could not resolve plugin dependencies", "error", err)

		return
	}

	for _, p := range ps {
		initPlugin(p)
	}
}

func UnloadPlugins() {
	plugin.CleanupClients()
}

func findPlugins() {
	matches, err := filepath.Glob(PathByName("*"))
	if err != nil {
		return
	}

	for _, path := range matches {
		logger := log.Default().With("path", path)

		if stat, err := os.Stat(path); err != nil || stat.Mode().Perm()&0111 == 0 {
			continue
		}

		for n, t := range pluginTypes {
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

			name := p.PluginInfo().Name

			if plugins[name.Type] == nil {
				plugins[name.Type] = map[string]Plugin{}
			}

			plugins[name.Type][name.Name] = p
		}
	}
}

func initPlugin(p Plugin) {
	info := p.PluginInfo()
	logger := log.L().With("name", info.Name.Name, "type", info.Name.Type)
	logger = augmentLogger(logger, info.Fields)

	if config, err := p.Init(); err != nil {
		logger.Warn("could not load plugin", "error", err)
		delete(plugins[info.Name.Type], info.Name.Name)
	} else {
		augmentLogger(logger, config).Info("loaded plugin")
	}
}

func loadGRPCPlugin(path string, pluginType string, grpcPlugin plugin.Plugin) (Plugin, error) {
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

	return nil, ErrPluginInvalid
}

func augmentLogger(logger log.Logger, fields map[string]string) log.Logger {
	fs := []interface{}{}

	for k, v := range fields {
		fs = append(fs, k, v)
	}

	return logger.With(fs...)
}
