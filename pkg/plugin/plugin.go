package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/oclaussen/dodo/pkg/appconfig"
)

const (
	ProtocolVersion  = 1
	MagicCookieKey   = "DODO_PLUGIN"
	MagicCookieValue = "69318785-d741-4150-ac91-8f03fa703530"

	ErrPluginInvalid      PluginError = "invalid plugin"
	ErrPluginNotFound     PluginError = "plugin not found"
	ErrNoValidPluginFound PluginError = "no valid plugin found"
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
	GRPCClient() plugin.Plugin
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

func LoadPlugins(types ...Type) {
	matches, err := filepath.Glob(fmt.Sprintf("%s/dodo-*_%s_%s", appconfig.GetPluginDir(), runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return
	}

	for _, path := range matches {
		logger := log.Default().With("path", path)

		if stat, err := os.Stat(path); err != nil || stat.Mode().Perm()&0111 == 0 {
			continue
		}

		pluginMap := map[string]plugin.Plugin{}
		for _, t := range types {
			pluginMap[t.String()] = t.GRPCClient()
		}

		client := plugin.NewClient(&plugin.ClientConfig{
			Managed:          true,
			Plugins:          pluginMap,
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
			logger.Debug("error getting plugin client", "error", err)
			continue
		}

		for _, t := range types {
			raw, err := conn.Dispense(t.String())
			if err != nil {
				logger.Debug("error dispensing plugin", "error", err)
				continue
			}

			p, ok := raw.(Plugin)
			if !ok {
				logger.Debug("plugin does not implement init")
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

func UnloadPlugins() {
	plugin.CleanupClients()
}
