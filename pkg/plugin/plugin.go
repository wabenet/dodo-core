package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-plugin"
	"github.com/oclaussen/dodo/pkg/appconfig"
	log "github.com/sirupsen/logrus"
)

type Plugin interface {
	Init() error
}

const (
	ProtocolVersion  = 1
	MagicCookieKey   = "DODO_PLUGIN"
	MagicCookieValue = "69318785-d741-4150-ac91-8f03fa703530"
)

var (
	serverPluginMap = map[string]plugin.Plugin{}
	clientPluginMap = map[string]plugin.Plugin{}
	plugins         = map[string][]Plugin{}
)

func RegisterPluginServer(t string, p plugin.Plugin) {
	serverPluginMap[t] = p
}

func RegisterPluginClient(t string, p plugin.Plugin) {
	clientPluginMap[t] = p
}

func ServePlugins() {
	log.SetFormatter(new(log.JSONFormatter))
	plugin.Serve(&plugin.ServeConfig{
		GRPCServer: plugin.DefaultGRPCServer,
		Plugins:    serverPluginMap,
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  ProtocolVersion,
			MagicCookieKey:   MagicCookieKey,
			MagicCookieValue: MagicCookieValue,
		},
	})
}

func GetPlugins(pluginType string) []Plugin {
	return plugins[pluginType]
}

func LoadPlugins() {
	matches, err := filepath.Glob(fmt.Sprintf("%s/dodo-*_%s_%s", appconfig.GetPluginDir(), runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return
	}

	for _, path := range matches {
		logger := log.WithFields(log.Fields{"path": path})

		if stat, err := os.Stat(path); err != nil || stat.Mode().Perm()&0111 == 0 {
			continue
		}

		client := plugin.NewClient(&plugin.ClientConfig{
			Managed:          true,
			Plugins:          clientPluginMap,
			Cmd:              exec.Command(path),
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
			Logger:           NewLogger(),
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  ProtocolVersion,
				MagicCookieKey:   MagicCookieKey,
				MagicCookieValue: MagicCookieValue,
			},
		})

		conn, err := client.Client()
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Debug("error getting plugin client")
			continue
		}

		for pluginType := range clientPluginMap {
			raw, err := conn.Dispense(pluginType)
			if err != nil {
				logger.WithFields(log.Fields{"error": err}).Debug("error dispensing plugin")
				continue
			}

			p, ok := raw.(Plugin)
			if !ok {
				logger.Debug("plugin does not implement init")
				continue
			}

			if err := p.Init(); err != nil {
				logger.WithFields(log.Fields{"error": err}).Debug("error initializing plugin")
				continue
			}

			logger.WithFields(log.Fields{"type": pluginType}).Debug("initialized plugin")

			plugins[pluginType] = append(plugins[pluginType], p)
		}
	}
}

func UnloadPlugins() {
	plugin.CleanupClients()
}
