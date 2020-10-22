package logx

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tOnkowzl/libs/contextx"
)

type Severity string

const (
	limitMSG = 5000

	severityKey = "severity"

	LevelInfo      Severity = "INFO"
	LevelDebug     Severity = "DEBUG"
	LevelWarn      Severity = "WARNING"
	LevelError     Severity = "ERROR"
	LevelEmergency Severity = "EMERGENCY"
)

func StandardLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

func Init(level, env string) {
	logrus.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(lvl)
	}

	if env == "dev" {
		logrus.SetFormatter(&logrus.TextFormatter{})
		return
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func LimitMSG(b []byte) string {
	if limitMSG < len(b) {
		return string(b[:limitMSG]) + "..."
	}
	return string(b)
}

func WithContext(ctx context.Context) logrus.FieldLogger {
	return logrus.WithField("id", contextx.GetID(ctx))
}

func WithSeverity(ctx context.Context, level Severity) logrus.FieldLogger {
	return logrus.WithFields(logrus.Fields{
		"id":        contextx.GetID(ctx),
		severityKey: level,
	})
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

func Debugf(format string, args ...interface{}) {
	logrus.WithField(severityKey, LevelDebug).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logrus.WithField(severityKey, LevelInfo).Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	logrus.WithField(severityKey, LevelInfo).Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.WithField(severityKey, LevelWarn).Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logrus.WithField(severityKey, LevelWarn).Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.WithField(severityKey, LevelError).Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.WithField(severityKey, LevelEmergency).Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.WithField(severityKey, LevelEmergency).Panicf(format, args...)
}

func Debug(args ...interface{}) {
	logrus.WithField(severityKey, LevelDebug).Debug(args...)
}

func Info(args ...interface{}) {
	logrus.WithField(severityKey, LevelInfo).Info(args...)
}

func Print(args ...interface{}) {
	logrus.WithField(severityKey, LevelInfo).Print(args...)
}

func Warn(args ...interface{}) {
	logrus.WithField(severityKey, LevelWarn).Warn(args...)
}

func Warning(args ...interface{}) {
	logrus.WithField(severityKey, LevelWarn).Warning(args...)
}

func Error(args ...interface{}) {
	logrus.WithField(severityKey, LevelError).Error(args...)
}

func Fatal(args ...interface{}) {
	logrus.WithField(severityKey, LevelEmergency).Fatal(args...)
}

func Panic(args ...interface{}) {
	logrus.WithField(severityKey, LevelEmergency).Panic(args...)
}

func Debugln(args ...interface{}) {
	logrus.WithField(severityKey, LevelDebug).Debugln(args...)
}

func Infoln(args ...interface{}) {
	logrus.WithField(severityKey, LevelInfo).Infoln(args...)
}

func Println(args ...interface{}) {
	logrus.WithField(severityKey, LevelInfo).Println(args...)
}

func Warnln(args ...interface{}) {
	logrus.WithField(severityKey, LevelWarn).Warnln(args...)
}

func Warningln(args ...interface{}) {
	logrus.WithField(severityKey, LevelWarn).Warningln(args...)
}

func Errorln(args ...interface{}) {
	logrus.WithField(severityKey, LevelError).Errorln(args...)
}

func Fatalln(args ...interface{}) {
	logrus.WithField(severityKey, LevelEmergency).Fatalln(args...)
}

func Panicln(args ...interface{}) {
	logrus.WithField(severityKey, LevelEmergency).Panicln(args...)
}
