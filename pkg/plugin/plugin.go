package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-plugin"
	"github.com/oclaussen/dodo/pkg/plugin/command"
	"github.com/oclaussen/dodo/pkg/plugin/configuration"
	"github.com/oclaussen/go-gimme/configfiles"
	log "github.com/sirupsen/logrus"
)

const (
	Command       = "command"
	Configuration = "configuration"
)

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "DODO_PLUGIN",
	MagicCookieValue: "plugin",
}

var PluginMap = map[string]plugin.Plugin{
	Command:       &command.Plugin{},
	Configuration: &configuration.Plugin{},
}

var clients []pluginClient

type pluginClient struct {
	path       string
	connection plugin.ClientProtocol
	cleanup    func()
}

func InitPlugins() {
	plugins := findPluginExecutables()
	for _, path := range plugins {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig:  HandshakeConfig,
			Plugins:          PluginMap,
			Cmd:              exec.Command(path),
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
			Logger:           NewPluginLogger(),
		})
		conn, err := client.Client()
		if err != nil {
			continue
		}
		clients = append(clients, pluginClient{
			path:       path,
			connection: conn,
			cleanup:    client.Kill,
		})
	}
	return
}

func FreePlugins() {
	for _, client := range clients {
		client.cleanup()
	}
}

func GetCommands() []command.Command {
	result := []command.Command{}
	for _, client := range clients {
		if raw, err := client.connection.Dispense(Command); err == nil {
			log.WithFields(log.Fields{"path": client.path, "type": Command}).Debug("found plugin")
			result = append(result, raw.(command.Command))
		}
	}
	return result
}

func GetConfigurations() []configuration.Configuration {
	result := []configuration.Configuration{}
	for _, client := range clients {
		if raw, err := client.connection.Dispense(Configuration); err == nil {
			log.WithFields(log.Fields{"path": client.path, "type": Configuration}).Debug("found plugin")
			result = append(result, raw.(configuration.Configuration))
		}
	}
	return result
}

func findPluginExecutables() []string {
	result := []string{}
	directories, err := configfiles.GimmeConfigDirectories(&configfiles.Options{Name: "dodo"})
	if err != nil {
		return result
	}
	for _, dir := range directories {
		matches, err := filepath.Glob(fmt.Sprintf("%s/plugins/dodo-*_%s_%s", dir, runtime.GOOS, runtime.GOARCH))
		if err != nil {
			continue
		}
		for _, path := range matches {
			if stat, err := os.Stat(path); err == nil && stat.Mode().Perm()&0111 != 0 {
				result = append(result, path)
			}
		}
	}
	return result
}
