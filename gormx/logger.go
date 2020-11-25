package gormx

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/tOnkowzl/libs/logx"
	glogger "gorm.io/gorm/logger"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

type Config struct {
	SlowThreshold time.Duration
	Colorful      bool
	LogLevel      glogger.LogLevel
}

var (
	Discard           = New(log.New(ioutil.Discard, "", log.LstdFlags), Config{})
	DefaultGormLogger = New(log.New(os.Stdout, "\r\n", log.LstdFlags), Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      glogger.Info,
		Colorful:      false,
	})
	Recorder = traceRecorder{Interface: DefaultGormLogger}
)

func New(writer glogger.Writer, config Config) glogger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &logger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type logger struct {
	glogger.Writer
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *logger) LogMode(level glogger.LogLevel) glogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= glogger.Info {
		logx.WithSeverityInfo(ctx).Infof(msg, data...)
	}
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= glogger.Warn {
		logx.WithSeverityWarn(ctx).Warnf(msg, data...)
	}
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= glogger.Error {
		logx.WithSeverityError(ctx).Errorf(msg, data...)
	}
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		sql, rows := fc()
		sql = logx.LimitMSGString(sql)
		switch {
		case err != nil && l.LogLevel >= glogger.Error:
			if rows == -1 {
				logx.WithSeverityError(ctx).Printf(l.traceErrStr, "", err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				logx.WithSeverityError(ctx).Printf(l.traceErrStr, "", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= glogger.Warn:
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				logx.WithSeverityWarn(ctx).Printf(l.traceWarnStr, "", slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				logx.WithSeverityWarn(ctx).Printf(l.traceWarnStr, "", slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel >= glogger.Info:
			if rows == -1 {
				logx.WithSeverityInfo(ctx).Printf(l.traceStr, "", float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				logx.WithSeverityInfo(ctx).Printf(l.traceStr, "", float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

type traceRecorder struct {
	glogger.Interface
	BeginAt      time.Time
	SQL          string
	RowsAffected int64
	Err          error
}

func (l traceRecorder) New() *traceRecorder {
	return &traceRecorder{Interface: l.Interface}
}

func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.BeginAt = begin
	l.SQL, l.RowsAffected = fc()
	l.Err = err
}
