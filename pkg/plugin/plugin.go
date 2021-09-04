package plugin

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/appconfig"
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

var (
	pluginTypes = map[string]Type{}
	plugins     = map[string]map[string]Plugin{}
)

type PluginError string

func (e PluginError) Error() string {
	return string(e)
}

type Plugin interface {
	Type() Type
	PluginInfo() (*api.PluginInfo, error)
}

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
		info, err := p.PluginInfo()
		if err != nil {
			log.L().Debug("plugin does not provide plugin info", "error", err)
			continue
		}

		t := p.Type().String()

		if plugins[t] == nil {
			plugins[t] = map[string]Plugin{}
		}

		plugins[p.Type().String()][info.Name] = p

		log.L().Debug("loaded plugin", "type", t, "name", info.Name)
	}
}

func ServePlugins(plugins ...Plugin) error {
	log.SetDefault(log.New(appconfig.GetPluginLoggerOptions()))

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

func GetPlugins(pluginType string) map[string]Plugin {
	return plugins[pluginType]
}

func PathByName(name string) string {
	return filepath.Join(
		appconfig.GetPluginDir(),
		fmt.Sprintf("dodo-%s_%s_%s", name, runtime.GOOS, runtime.GOARCH),
	)
}

func LoadPlugins() {
	matches, err := filepath.Glob(PathByName("*"))
	if err != nil {
		return
	}

	for _, path := range matches {
		logger := log.Default().With("path", path)

		if stat, err := os.Stat(path); err != nil || stat.Mode().Perm()&0111 == 0 {
			continue
		}

		for _, t := range pluginTypes {
			logger.Debug("attempt loading plugin", "type", t.String())

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

			IncludePlugins(p)
		}
	}
}

func loadGRPCPlugin(path string, pluginType string, grpcPlugin plugin.Plugin) (Plugin, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		Managed:          true,
		Plugins:          map[string]plugin.Plugin{pluginType: grpcPlugin},
		Cmd:              exec.Command(path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		Logger:           log.New(appconfig.GetLoggerOptions()),
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

func UnloadPlugins() {
	plugin.CleanupClients()
}
