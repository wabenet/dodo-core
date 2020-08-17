package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/dodo/dodo-core/pkg/appconfig"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

const (
	ProtocolVersion  = 1
	MagicCookieKey   = "DODO_PLUGIN"
	MagicCookieValue = "69318785-d741-4150-ac91-8f03fa703530"

	ErrPluginInvalid        PluginError = "invalid plugin"
	ErrPluginNotImplemented PluginError = "not implemented"
	ErrPluginNotFound       PluginError = "plugin not found"
	ErrNoValidPluginFound   PluginError = "no valid plugin found"
)

var plugins = map[string][]Plugin{}

type PluginError string

func (e PluginError) Error() string {
	return string(e)
}

type Plugin interface {
	Type() Type
	Init() error
}

type Type interface {
	String() string
	GRPCClient() (plugin.Plugin, error)
	GRPCServer(Plugin) (plugin.Plugin, error)
}

func RegisterBuiltin(p Plugin) {
	if err := p.Init(); err != nil {
		log.L().Debug("error initializing plugin", "error", err)
		return
	}
	plugins[p.Type().String()] = append(plugins[p.Type().String()], p)
}

func ServePlugins(plugins ...Plugin) error {
	pluginMap := map[string]plugin.Plugin{}
	for _, p := range plugins {
		s, err := p.Type().GRPCServer(p)
		if err != nil {
			return err
		}

		pluginMap[p.Type().String()] = s
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

func GetPlugins(pluginType string) []Plugin {
	return plugins[pluginType]
}

func PathByName(name string) string {
	return filepath.Join(
		appconfig.GetPluginDir(),
		fmt.Sprintf("dodo-%s_%s_%s", name, runtime.GOOS, runtime.GOARCH),
	)
}

func LoadPlugins(types ...Type) {
	matches, err := filepath.Glob(PathByName("*"))
	if err != nil {
		return
	}

	for _, path := range matches {
		logger := log.Default().With("path", path)

		if stat, err := os.Stat(path); err != nil || stat.Mode().Perm()&0111 == 0 {
			continue
		}

		for _, t := range types {
			grpcClient, err := t.GRPCClient()
			if err != nil {
				logger.Debug("error loading plugin", "error", err)
				continue
			}

			p, err := loadGRPCPlugin(path, t.String(), grpcClient)
			if err != nil {
				logger.Debug("could not load plugin over grpc", "error", err)
				continue
			}

			if err := p.Init(); err != nil {
				logger.Debug("error initializing plugin", "error", err)
				continue
			}

			logger.Debug("initialized plugin", "type", t.String())

			plugins[t.String()] = append(plugins[t.String()], p)
		}
	}
}

func loadGRPCPlugin(path string, pluginType string, grpcPlugin plugin.Plugin) (Plugin, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		Managed:          true,
		Plugins:          map[string]plugin.Plugin{pluginType: grpcPlugin},
		Cmd:              exec.Command(path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		Logger:           log.Default(),
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  ProtocolVersion,
			MagicCookieKey:   MagicCookieKey,
			MagicCookieValue: MagicCookieValue,
		},
	})

	conn, err := client.Client()
	if err != nil {
		return nil, fmt.Errorf("error getting plugin client: %w", err)
	}

	raw, err := conn.Dispense(pluginType)
	if err != nil {
		return nil, fmt.Errorf("error dispensing plugin: %w", err)
	}

	p, ok := raw.(Plugin)
	if !ok {
		return nil, ErrPluginInvalid
	}
	return p, nil
}

func UnloadPlugins() {
	plugin.CleanupClients()
}
