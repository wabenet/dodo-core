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

const (
	MagicCookieKey   = "DODO_PLUGIN"
	MagicCookieValue = "69318785-d741-4150-ac91-8f03fa703530"
)

var (
	serverPluginMap = map[string]plugin.Plugin{}
	clientPluginMap = map[string]plugin.Plugin{}
	clients         = []plugin.ClientProtocol{}
	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   MagicCookieKey,
		MagicCookieValue: MagicCookieValue,
	}
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
		GRPCServer:      plugin.DefaultGRPCServer,
		HandshakeConfig: HandshakeConfig,
		Plugins:         serverPluginMap,
	})
}

func GetPlugins(pluginType string) []interface{} {
	result := []interface{}{}
	for _, p := range clients {
		raw, err := p.Dispense(pluginType)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Debug("error dispensing plugin")
			continue
		} else {
			result = append(result, raw)
		}
	}
	return result
}

func LoadPlugins() {
	for _, path := range loadPluginExecutables() {
		client := plugin.NewClient(&plugin.ClientConfig{
			Managed:          true,
			HandshakeConfig:  HandshakeConfig,
			Plugins:          clientPluginMap,
			Cmd:              exec.Command(path),
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
			Logger:           NewPluginLogger(),
		})

		conn, err := client.Client()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Debug("error getting plugin client")
		} else {
			log.WithFields(log.Fields{"path": path}).Debug("found plugin")
			clients = append(clients, conn)
		}
	}
}

func UnloadPlugins() {
	plugin.CleanupClients()
}

func loadPluginExecutables() []string {
	executables := []string{}
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
	return executables
}
