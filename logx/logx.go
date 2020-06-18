package logx

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	limitMSG = 5000
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

func WithID(requestID string) logrus.FieldLogger {
	return logrus.WithField("id", requestID)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	logrus.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Print(args ...interface{}) {
	logrus.Print(args...)
}

func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

func Warning(args ...interface{}) {
	logrus.Warning(args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

func Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

func Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

func Println(args ...interface{}) {
	logrus.Println(args...)
}

func Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}

func Warningln(args ...interface{}) {
	logrus.Warningln(args...)
}

func Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	logrus.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	logrus.Panicln(args...)
}
