package logx

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tOnkowzl/libs/contextx"
)

type Severity string

const (
	limitMSG = 2500

	severityKey = "severity"

	SeverityInfo      Severity = "INFO"
	SeverityDebug     Severity = "DEBUG"
	SeverityWarn      Severity = "WARNING"
	SeverityError     Severity = "ERROR"
	SeverityEmergency Severity = "EMERGENCY"
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

func LimitMSGByte(b []byte) string {
	if limitMSG < len(b) {
		return string(b[:limitMSG]) + "..."
	}
	return string(b)
}

func LimitMSGString(s string) string {
	if limitMSG < len(s) {
		return s[:limitMSG] + "..."
	}
	return s
}

func WithContext(ctx context.Context) logrus.FieldLogger {
	return logrus.WithField("id", contextx.GetID(ctx))
}

func WithSeverityInfo(ctx context.Context) logrus.FieldLogger {
	return logrus.WithFields(logrus.Fields{
		"id":        contextx.GetID(ctx),
		severityKey: SeverityInfo,
	})
}

func WithSeverityDebug(ctx context.Context) logrus.FieldLogger {
	return logrus.WithFields(logrus.Fields{
		"id":        contextx.GetID(ctx),
		severityKey: SeverityDebug,
	})
}

func WithSeverityWarn(ctx context.Context) logrus.FieldLogger {
	return logrus.WithFields(logrus.Fields{
		"id":        contextx.GetID(ctx),
		severityKey: SeverityWarn,
	})
}

func WithSeverityError(ctx context.Context) logrus.FieldLogger {
	return logrus.WithFields(logrus.Fields{
		"id":        contextx.GetID(ctx),
		severityKey: SeverityError,
	})
}

func WithSeverityEmergency(ctx context.Context) logrus.FieldLogger {
	return logrus.WithFields(logrus.Fields{
		"id":        contextx.GetID(ctx),
		severityKey: SeverityEmergency,
	})
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

func Debugf(format string, args ...interface{}) {
	logrus.WithField(severityKey, SeverityDebug).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logrus.WithField(severityKey, SeverityInfo).Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	logrus.WithField(severityKey, SeverityInfo).Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.WithField(severityKey, SeverityWarn).Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logrus.WithField(severityKey, SeverityWarn).Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.WithField(severityKey, SeverityError).Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.WithField(severityKey, SeverityEmergency).Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.WithField(severityKey, SeverityEmergency).Panicf(format, args...)
}

func Debug(args ...interface{}) {
	logrus.WithField(severityKey, SeverityDebug).Debug(args...)
}

func Info(args ...interface{}) {
	logrus.WithField(severityKey, SeverityInfo).Info(args...)
}

func Print(args ...interface{}) {
	logrus.WithField(severityKey, SeverityInfo).Print(args...)
}

func Warn(args ...interface{}) {
	logrus.WithField(severityKey, SeverityWarn).Warn(args...)
}

func Warning(args ...interface{}) {
	logrus.WithField(severityKey, SeverityWarn).Warning(args...)
}

func Error(args ...interface{}) {
	logrus.WithField(severityKey, SeverityError).Error(args...)
}

func Fatal(args ...interface{}) {
	logrus.WithField(severityKey, SeverityEmergency).Fatal(args...)
}

func Panic(args ...interface{}) {
	logrus.WithField(severityKey, SeverityEmergency).Panic(args...)
}

func Debugln(args ...interface{}) {
	logrus.WithField(severityKey, SeverityDebug).Debugln(args...)
}

func Infoln(args ...interface{}) {
	logrus.WithField(severityKey, SeverityInfo).Infoln(args...)
}

func Println(args ...interface{}) {
	logrus.WithField(severityKey, SeverityInfo).Println(args...)
}

func Warnln(args ...interface{}) {
	logrus.WithField(severityKey, SeverityWarn).Warnln(args...)
}

func Warningln(args ...interface{}) {
	logrus.WithField(severityKey, SeverityWarn).Warningln(args...)
}

func Errorln(args ...interface{}) {
	logrus.WithField(severityKey, SeverityError).Errorln(args...)
}

func Fatalln(args ...interface{}) {
	logrus.WithField(severityKey, SeverityEmergency).Fatalln(args...)
}

func Panicln(args ...interface{}) {
	logrus.WithField(severityKey, SeverityEmergency).Panicln(args...)
}
