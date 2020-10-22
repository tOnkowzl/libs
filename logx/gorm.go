package logx

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
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

var (
	Discard           = New(log.New(ioutil.Discard, "", log.LstdFlags), glogger.Config{})
	DefaultGormLogger = New(log.New(os.Stdout, "\r\n", log.LstdFlags), glogger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      glogger.Info,
		Colorful:      true,
	})
	Recorder = traceRecorder{Interface: DefaultGormLogger}
)

func New(writer glogger.Writer, config glogger.Config) glogger.Interface {
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
	glogger.Config
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
		WithContext(ctx).Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= glogger.Warn {
		WithContext(ctx).Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= glogger.Error {
		WithContext(ctx).Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= glogger.Error:
			sql, rows := fc()
			if rows == -1 {
				WithContext(ctx).Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				WithContext(ctx).Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= glogger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				WithContext(ctx).Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				WithContext(ctx).Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel >= glogger.Info:
			sql, rows := fc()
			if rows == -1 {
				WithContext(ctx).Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				WithContext(ctx).Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
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
