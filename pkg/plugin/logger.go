package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/sirupsen/logrus"
)

type PluginLogger struct {
	name   string
	logger logrus.FieldLogger
}

func NewPluginLogger() hclog.Logger {
	return &PluginLogger{
		logger: &logrus.Logger{
			Out:   os.Stderr,
			Level: logrus.GetLevel(),
			Formatter: &logrus.TextFormatter{
				DisableTimestamp:       true,
				DisableLevelTruncation: true,
			},
		},
	}
}

func (logger *PluginLogger) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.Debug:
		logger.Debug(msg, args...)
	case hclog.Info:
		logger.Info(msg, args...)
	case hclog.Warn:
		logger.Warn(msg, args...)
	case hclog.Error:
		logger.Error(msg, args...)
	}
}

func (*PluginLogger) Trace(_ string, _ ...interface{}) {
}

func (logger *PluginLogger) IsTrace() bool {
	return false
}

func (logger *PluginLogger) Debug(msg string, args ...interface{}) {
	// All plugin output logs to debug
	done := logger.upgradePluginOutput(msg, args)
	if !done {
		logger.logger.WithFields(argsToFields(args)).Debug(msg)
	}
}

func (logger *PluginLogger) IsDebug() bool {
	return logger.logger.WithFields(logrus.Fields{}).Level >= logrus.DebugLevel
}

func (logger *PluginLogger) Info(msg string, args ...interface{}) {
	logger.logger.WithFields(argsToFields(args)).Info(msg)
}

func (logger *PluginLogger) IsInfo() bool {
	return logger.logger.WithFields(logrus.Fields{}).Level >= logrus.InfoLevel
}

func (logger *PluginLogger) Warn(msg string, args ...interface{}) {
	logger.logger.WithFields(argsToFields(args)).Warn(msg)
}

func (logger *PluginLogger) IsWarn() bool {
	return logger.logger.WithFields(logrus.Fields{}).Level >= logrus.WarnLevel
}

func (logger *PluginLogger) Error(msg string, args ...interface{}) {
	logger.logger.WithFields(argsToFields(args)).Error(msg)
}

func (logger *PluginLogger) IsError() bool {
	return logger.logger.WithFields(logrus.Fields{}).Level >= logrus.ErrorLevel
}

func (logger *PluginLogger) SetLevel(_ hclog.Level) {}

func (logger *PluginLogger) With(args ...interface{}) hclog.Logger {
	return &PluginLogger{logger: logger.logger.WithFields(argsToFields(args))}
}

func (logger *PluginLogger) ImpliedArgs() []interface{} {
	return []interface{}{}
}

func (logger *PluginLogger) Name() string {
	return logger.name
}

func (logger *PluginLogger) Named(name string) hclog.Logger {
	if len(logger.name) > 0 {
		return logger.ResetNamed(fmt.Sprintf("%s.%s", logger.name, name))
	}
	return logger.ResetNamed(name)
}

func (logger *PluginLogger) ResetNamed(name string) hclog.Logger {
	return &PluginLogger{name: name, logger: logger.logger}
}

func (logger *PluginLogger) StandardLogger(_ *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(logger.logger.WithFields(logrus.Fields{}).WriterLevel(logrus.InfoLevel), "", 0)
}

func (logger *PluginLogger) StandardWriter(_ *hclog.StandardLoggerOptions) io.Writer {
	if l, ok := logger.logger.(*logrus.Logger); ok {
		return l.Out
	}
	return os.Stderr
}

func (logger *PluginLogger) upgradePluginOutput(originalMsg string, args []interface{}) bool {
	var output map[string]string
	if err := json.Unmarshal([]byte(originalMsg), &output); err != nil {
		return false
	}

	var msg, level string
	fields := argsToFields(args)
	for k, v := range output {
		switch k {
		case "msg":
			msg = v
		case "level":
			level = v
		case "time":
			// Time will be overridden by parent logger
		default:
			fields[k] = v
		}
	}

	switch level {
	case "trace":
		logger.logger.WithFields(fields).Trace(msg)
	case "debug":
		logger.logger.WithFields(fields).Debug(msg)
	case "info":
		logger.logger.WithFields(fields).Info(msg)
	case "warn", "warning":
		logger.logger.WithFields(fields).Warn(msg)
	case "error":
		logger.logger.WithFields(fields).Error(msg)
	case "fatal":
		logger.logger.WithFields(fields).Fatal(msg)
	case "panic":
		logger.logger.WithFields(fields).Panic(msg)
	default:
		return false
	}
	return true
}

func argsToFields(args []interface{}) logrus.Fields {
	if len(args)%2 != 0 {
		args = append(args, "")
	}
	fields := make(logrus.Fields, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		if key, ok := args[i].(string); ok {
			fields[key] = args[i+1]
		}
	}
	return fields
}
