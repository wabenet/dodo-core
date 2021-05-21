package appconfig

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"syscall"

	log "github.com/hashicorp/go-hclog"
)

const (
	systemAppDir = "/var/lib/dodo"
)

func GetAppDir() string {
	dir := filepath.FromSlash(systemAppDir)
	if user, err := user.Current(); err == nil && user.HomeDir != "" {
		dir = filepath.Join(user.HomeDir, ".dodo")
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		log.L().Warn("could not create app directory", "error", err)
	}

	return dir
}

func GetPluginDir() string {
	return filepath.Join(GetAppDir(), "plugins")
}

func GetLoggerOptions() *log.LoggerOptions {
	output := ioutil.Discard
	level := log.Info

	if levelName := os.Getenv("DODO_LOG_LEVEL"); levelName != "" {
		switch levelName {
		case "ERROR":
			level = log.Error
		case "WARN":
			level = log.Warn
		case "INFO":
			level = log.Info
		case "DEBUG":
			level = log.Debug
		case "TRACE":
			level = log.Trace
		}
	}

	if path := os.Getenv("DODO_LOG_PATH"); path != "" {
		if outFile, err := os.OpenFile(path, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666); err == nil {
			output = outFile
		}
	}

	output = os.Stderr

	return &log.LoggerOptions{
		Name:       "dodo",
		Level:      level,
		Output:     output,
		JSONFormat: true,
	}
}

func GetPluginLoggerOptions() *log.LoggerOptions {
	level := log.Info

	if levelName := os.Getenv("DODO_LOG_LEVEL"); levelName != "" {
		switch levelName {
		case "ERROR":
			level = log.Error
		case "WARN":
			level = log.Warn
		case "INFO":
			level = log.Info
		case "DEBUG":
			level = log.Debug
		case "TRACE":
			level = log.Trace
		}
	}

	return &log.LoggerOptions{
		Name:       "dodo",
		Level:      level,
		Output:     os.Stderr,
		JSONFormat: true,
	}
}
