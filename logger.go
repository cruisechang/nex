package nex

import (
	"time"
	"os"
	"fmt"
	golog "log"
	"runtime"
)

const (
	LogLevelDebug   uint8 = 0
	LogLevelInfo    uint8 = 1
	LogLevelWarning uint8 = 2
	LogLevelError   uint8 = 3
	LogLevelPanic   uint8 = 4
)

type Logger interface {
	SetLogFileName(name string)
	SetLevel(level uint8)
	ActiveLog(active bool)
	ActiveLogFile(acvite bool)
	Log(logLevel uint8, msg string)
	LogFile(logLevel uint8, v ...interface{})
}

type logger struct {
	level         uint8
	logFileName   string
	activeLog     bool
	activeLogFile bool
}

func NewLogger() (Logger, error) {
	l := &logger{
		level:         LogLevelWarning,
		logFileName:   "nex",
		activeLog:     true,
		activeLogFile: true,
	}
	return l, nil
}
func (l *logger) SetLogFileName(name string) {
	l.logFileName = name
}
func (l *logger) SetLevel(level uint8) {
	l.level = level
}
func (l *logger) ActiveLog(active bool) {
	l.activeLog = active
}
func (l *logger) ActiveLogFile(active bool) {
	l.activeLogFile = active
}
func (l *logger) Log(logLevel uint8, msg string) {

	if !l.activeLog {
		return
	}

	prefix := "[Info]"
	switch logLevel {
	case LogLevelDebug:
		prefix = "[Debug]"
	case LogLevelWarning:
		prefix = "[warning]"
	case LogLevelError:
		prefix = "[Error]"
	case LogLevelPanic:
		prefix = "[Panic]"

	}

	_, file, line, _ := runtime.Caller(1)

	golog.Printf("%s%s:%d %v", prefix, file, line, msg)
}

func (l *logger) LogFile(logLevel uint8, v ...interface{}) {

	if !l.activeLogFile {
		return
	}

	if logLevel >= l.level {

		prefix := "[Info]"
		switch logLevel {
		case LogLevelDebug:
			prefix = "[Debug]"
		case LogLevelWarning:
			prefix = "[warning]"
		case LogLevelError:
			prefix = "[Error]"
		case LogLevelPanic:
			prefix = "[Panic]"

		}
		filePath := l.logFileName + time.Now().Format("20060102") + ".log"
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err == nil {
			lr := golog.New(f, prefix, golog.LstdFlags|golog.Lshortfile)
			lr.Output(5, fmt.Sprintln(v))

		}
		defer f.Close()

	}
}
