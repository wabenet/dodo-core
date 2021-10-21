package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

const (
	Name = "dodo"

	ConfKeyLogLevel = "log-level"
	ConfKeyLogFile  = "log-file"
	ConfKeyAppDir   = "app-dir"

	DefaultLogLevel = "INFO"
	DefaultAppDir   = "/var/lib/dodo"
)

func Configure() {
	viper.SetConfigName(Name)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix(Name)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.AddConfigPath(fmt.Sprintf("/etc/%s", Name))

	viper.SetDefault(ConfKeyLogLevel, DefaultLogLevel)

	if user, err := user.Current(); err == nil && user.HomeDir != "" {
		dotDir := filepath.Join(user.HomeDir, fmt.Sprintf(".%s", Name))

		viper.AddConfigPath(dotDir)
		viper.SetDefault(ConfKeyAppDir, dotDir)
	} else {
		viper.SetDefault(ConfKeyAppDir, DefaultAppDir)
	}

	// Init logger with the default values, so the following log lines
	// already use our logging config.
	log.SetDefault(log.New(GetLoggerOptions()))

	if err := viper.ReadInConfig(); err != nil {
		var e *viper.ConfigFileNotFoundError
		if errors.As(err, &e) {
			log.L().Warn("no configuration file found", "error", err)
		} else {
			log.L().Warn("could not read config file", "error", err)
		}
	}
}

func GetAppDir() string {
	dir := filepath.FromSlash(viper.GetString(ConfKeyAppDir))

	if err := os.MkdirAll(dir, 0700); err != nil {
		log.L().Warn("could not create app directory", "error", err)
	}

	return dir
}

func GetPluginDir() string {
	return filepath.Join(GetAppDir(), "plugins")
}

func LogLevel() log.Level {
	switch viper.GetString(ConfKeyLogLevel) {
	case "ERROR":
		return log.Error
	case "WARN":
		return log.Warn
	case "INFO":
		return log.Info
	case "DEBUG":
		return log.Debug
	case "TRACE":
		return log.Trace
	default:
		return log.Info
	}
}

func LogOutput() io.Writer {
	switch path := viper.GetString(ConfKeyLogFile); path {
	case "":
		return io.Discard
	case "-":
		return os.Stderr
	default:
		if outFile, err := os.OpenFile(path, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666); err == nil {
			return outFile
		}
		return io.Discard
	}
}

func GetLoggerOptions() *log.LoggerOptions {
	return &log.LoggerOptions{
		Name:       Name,
		Level:      LogLevel(),
		Output:     LogOutput(),
		JSONFormat: true,
	}
}

func GetPluginLoggerOptions() *log.LoggerOptions {
	return &log.LoggerOptions{
		Name:       Name,
		Level:      LogLevel(),
		Output:     os.Stderr,
		JSONFormat: true,
	}
}
