package log

import (
	"fmt"
	goLog "log"
	"os"
	"runtime"
	"time"
)

type Level uint8

const (
	LevelDebug   Level = 0
	LevelInfo    Level = 1
	LevelWarning Level = 2
	LevelError   Level = 3
	LevelPanic   Level = 4
)

type Logger interface {
	SetLogFileName(name string)
	SetLevel(level Level)
	ActiveLog(active bool)
	ActiveLogFile(active bool)
	Log(level Level, msg string)
	LogFile(level Level, v ...interface{})
}

type logger struct {
	level         Level
	logFileName   string
	activeLog     bool
	activeLogFile bool
	file          *os.File
}

func NewLogger() (Logger, error) {
	l := &logger{
		level:         LevelWarning,
		logFileName:   "nex",
		activeLog:     true,
		activeLogFile: true,
	}

	f, err := l.createFile()
	if err != nil {
		return nil, err
	}
	l.file = f
	return l, nil
}
func (l *logger) SetLogFileName(name string) {
	l.logFileName = name
}
func (l *logger) SetLevel(level Level) {
	l.level = level
}
func (l *logger) ActiveLog(active bool) {
	l.activeLog = active
}
func (l *logger) ActiveLogFile(active bool) {
	l.activeLogFile = active
}
func (l *logger) Log(logLevel Level, msg string) {

	if !l.activeLog {
		return
	}

	prefix := "[Info]"
	switch logLevel {
	case LevelDebug:
		prefix = "[Debug]"
	case LevelWarning:
		prefix = "[warning]"
	case LevelError:
		prefix = "[Error]"
	case LevelPanic:
		prefix = "[Panic]"

	}

	_, file, line, _ := runtime.Caller(1)

	goLog.Printf("%s%s:%d %v", prefix, file, line, msg)
}
func (l *logger) createFile() (*os.File, error) {
	filePath := l.logFileName + time.Now().Format("20060102") + ".log"
	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

}
func (l *logger) closeFile() error {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			goLog.Printf("LogFile error=%s", err.Error())
			return err
		}
	}
	return nil
}

func (l *logger) LogFile(logLevel Level, v ...interface{}) {

	if !l.activeLogFile {
		return
	}

	if logLevel >= l.level {

		prefix := "[Info]"
		switch logLevel {
		case LevelDebug:
			prefix = "[Debug]"
		case LevelWarning:
			prefix = "[warning]"
		case LevelError:
			prefix = "[Error]"
		case LevelPanic:
			prefix = "[Panic]"

		}

		filePath := l.logFileName + time.Now().Format("20060102") + ".log"
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		defer func() {
			if err := f.Close(); err != nil {
				//panic(err)
				goLog.Printf("LogFile error=%s", err.Error())
			}
		}()

		if err == nil {
			lr := goLog.New(f, prefix, goLog.LstdFlags|goLog.Lshortfile)
			lr.Output(5, fmt.Sprintln(v))
		}

		//if l.file != nil {
		//	lr := goLog.New(l.file, prefix, goLog.LstdFlags|goLog.Lshortfile)
		//	lr.Output(5, fmt.Sprintln(v))
		//} else {
		//	file, err := l.createFile()
		//	if err != nil {
		//		goLog.Printf("LogFile file=nil, then createFile() err%s", err.Error())
		//		return
		//
		//	}
		//	l.file = file
		//	l.LogFile(logLevel, v)
		//}
	}
}
