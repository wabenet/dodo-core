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
	logger *logrus.Logger
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

func (l *PluginLogger) Trace(msg string, args ...interface{}) {
	l.Log(hclog.Trace, msg, args...)
}

func (l *PluginLogger) IsTrace() bool {
	return l.logger.Level <= logrus.TraceLevel
}

func (l *PluginLogger) Debug(msg string, args ...interface{}) {
	l.Log(hclog.Debug, msg, args...)
}

func (l *PluginLogger) IsDebug() bool {
	return l.logger.Level <= logrus.DebugLevel
}

func (l *PluginLogger) Info(msg string, args ...interface{}) {
	l.Log(hclog.Info, msg, args...)
}

func (l *PluginLogger) IsInfo() bool {
	return l.logger.Level <= logrus.InfoLevel
}

func (l *PluginLogger) Warn(msg string, args ...interface{}) {
	l.Log(hclog.Warn, msg, args...)
}

func (l *PluginLogger) IsWarn() bool {
	return l.logger.Level <= logrus.WarnLevel
}

func (l *PluginLogger) Error(msg string, args ...interface{}) {
	l.Log(hclog.Error, msg, args...)
}

func (l *PluginLogger) IsError() bool {
	return l.logger.Level <= logrus.ErrorLevel
}

func (l *PluginLogger) SetLevel(level hclog.Level) {
	switch level {
	case hclog.Trace:
		l.logger.SetLevel(logrus.TraceLevel)
	case hclog.Debug:
		l.logger.SetLevel(logrus.DebugLevel)
	case hclog.Info:
		l.logger.SetLevel(logrus.InfoLevel)
	case hclog.Warn:
		l.logger.SetLevel(logrus.WarnLevel)
	case hclog.Error:
		l.logger.SetLevel(logrus.ErrorLevel)
	}
}

func (l *PluginLogger) With(args ...interface{}) hclog.Logger {
	return l // TODO
}

func (l *PluginLogger) ImpliedArgs() []interface{} {
	return []interface{}{} // TODO
}

func (l *PluginLogger) Name() string {
	return l.name
}

func (l *PluginLogger) Named(name string) hclog.Logger {
	if len(l.name) > 0 {
		return l.ResetNamed(fmt.Sprintf("%s.%s", l.name, name))
	}

	return l.ResetNamed(name)
}

func (l *PluginLogger) ResetNamed(name string) hclog.Logger {
	return &PluginLogger{name: name, logger: l.logger}
}

func (l *PluginLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(l.StandardWriter(opts), "", 0)
}

func (l *PluginLogger) StandardWriter(_ *hclog.StandardLoggerOptions) io.Writer {
	return l.logger.Out
}

func (l *PluginLogger) Log(level hclog.Level, msg string, args ...interface{}) {
	fields := argsToFields(args)

	var output map[string]json.RawMessage
	if err := json.Unmarshal([]byte(msg), &output); err != nil {
		switch level {
		case hclog.Trace:
			l.logger.WithFields(fields).Trace(msg)
		case hclog.Debug:
			l.logger.WithFields(fields).Debug(msg)
		case hclog.Info:
			l.logger.WithFields(fields).Info(msg)
		case hclog.Warn:
			l.logger.WithFields(fields).Warn(msg)
		case hclog.Error:
			l.logger.WithFields(fields).Error(msg)
		}

		return
	}

	var newMsg, newLevel string

	for k, v := range output {
		switch k {
		case "msg":
			json.Unmarshal(v, &newMsg)
		case "level":
			json.Unmarshal(v, &newLevel)
		case "time":
			// Time will be overridden by parent logger
		default:
			var data interface{}

			json.Unmarshal(v, &data)
			fields[k] = data
		}
	}

	switch newLevel {
	case "trace":
		l.logger.WithFields(fields).Trace(newMsg)
	case "debug":
		l.logger.WithFields(fields).Debug(newMsg)
	case "info":
		l.logger.WithFields(fields).Info(newMsg)
	case "warn", "warning":
		l.logger.WithFields(fields).Warn(newMsg)
	case "error":
		l.logger.WithFields(fields).Error(newMsg)
	case "fatal":
		l.logger.WithFields(fields).Fatal(newMsg)
	case "panic":
		l.logger.WithFields(fields).Panic(newMsg)
	default:
		l.logger.WithFields(fields).Debug(newMsg)
	}
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
