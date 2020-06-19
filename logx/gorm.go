package logx

import (
	"time"

	"github.com/sirupsen/logrus"
)

type GormLogger struct{}

func NewGormLogger() *GormLogger {
	return new(GormLogger)
}

func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		logrus.WithFields(logrus.Fields{
			"rows_returned": v[5],
			"src":           v[1],
			"values":        v[4],
			"duration":      time.Duration(v[2].(int64)).String(),
			"stmt":          v[3],
		}).Info("gorm response information")
	case "log":
		logrus.WithFields(logrus.Fields{"data": v[2]}).Print("gorm log")
	}
}
