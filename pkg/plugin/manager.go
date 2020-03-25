package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-plugin"
	"github.com/oclaussen/go-gimme/configfiles"
	log "github.com/sirupsen/logrus"
)

var (
	serverPluginMap = map[string]plugin.Plugin{}
	clientPluginMap = map[string]plugin.Plugin{}
	executables     = []string{}
	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "DODO_PLUGIN",
		MagicCookieValue: "plugin",
	}
)

func RegisterPluginServer(t string, p plugin.Plugin) {
	serverPluginMap[t] = p
}

func RegisterPluginClient(t string, p plugin.Plugin) {
	clientPluginMap[t] = p
}

type PluginManager struct {
	Plugins map[string]interface{}
	clients []*plugin.Client
}

func ServePlugins() {
	log.SetFormatter(new(log.JSONFormatter))
	plugin.Serve(&plugin.ServeConfig{
		GRPCServer:      plugin.DefaultGRPCServer,
		HandshakeConfig: HandshakeConfig,
		Plugins:         serverPluginMap,
	})
}

func LoadPlugins(pluginType string) *PluginManager {
	if len(executables) == 0 {
		loadPluginExecutables()
	}

	result := &PluginManager{
		Plugins: map[string]interface{}{},
		clients: []*plugin.Client{},
	}

	for _, path := range executables {
		logger := log.WithFields(log.Fields{"path": path, "type": pluginType})

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig:  HandshakeConfig,
			Plugins:          clientPluginMap,
			Cmd:              exec.Command(path),
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
			Logger:           NewPluginLogger(),
		})

		conn, err := client.Client()
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Debug("error getting plugin client")
			client.Kill()
			continue
		}

		raw, err := conn.Dispense(pluginType)
		if err != nil {
			logger.WithFields(log.Fields{"error": err}).Debug("error dispensing plugin")
			client.Kill()
			continue
		}

		logger.Debug("found plugin")
		result.clients = append(result.clients, client)
		result.Plugins[path] = raw
	}

	return result
}

func (m *PluginManager) UnloadPlugins() {
	for _, client := range m.clients {
		log.Debug("killing client")
		client.Kill()
	}
}

func loadPluginExecutables() {
	if directories, err := configfiles.GimmeConfigDirectories(&configfiles.Options{Name: "dodo"}); err == nil {
		for _, dir := range directories {
			matches, err := filepath.Glob(fmt.Sprintf("%s/plugins/dodo-*_%s_%s", dir, runtime.GOOS, runtime.GOARCH))
			if err != nil {
				continue
			}
			for _, path := range matches {
				if stat, err := os.Stat(path); err == nil && stat.Mode().Perm()&0111 != 0 {
					executables = append(executables, path)
				}
			}
		}
	}
}
