package applog

import (
	"fmt"
	"io"
	"log"
	"os"
)

var Log = New(os.Stdout)

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Stackf(string, ...interface{})
}

type Level uint8

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
	STACK
)

var (
	levelFmt = "level:%v\t"
	LogLevel = DEBUG
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	case STACK:
		return "STACK"
	default:
		return "<unknown>"
	}
}

const (
	standardLogFlags = log.Ltime | log.Ldate | log.Lmicroseconds
)

type logger struct {
	logger *log.Logger
}

func New(w io.Writer) *logger {
	newLogger := log.New(w, "", standardLogFlags)
	return &logger{logger: newLogger}
}

func shouldLog(level Level) bool {
	return level >= LogLevel
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if shouldLog(DEBUG) {
		lvl := fmt.Sprintf(levelFmt, DEBUG)
		msg := fmt.Sprintf(format, v...)
		lvlWithMsg := lvl + msg
		l.logger.Print(lvlWithMsg)
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if shouldLog(INFO) {
		lvl := fmt.Sprintf(levelFmt, INFO)
		msg := fmt.Sprintf(format, v...)
		l.logger.Print(lvl + msg)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if shouldLog(WARN) {
		lvl := fmt.Sprintf(levelFmt, WARN)
		msg := fmt.Sprintf(format, v...)
		l.logger.Print(lvl + msg)
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if shouldLog(ERROR) {
		lvl := fmt.Sprintf(levelFmt, ERROR)
		msg := fmt.Sprintf(format, v...)
		l.logger.Print(lvl + msg)
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	if shouldLog(FATAL) {
		lvl := fmt.Sprintf(levelFmt, FATAL)
		msg := fmt.Sprintf(format, v...)
		l.logger.Fatal(lvl + msg)
	}
}

func (l *logger) Stackf(format string, v ...interface{}) {
	lvl := fmt.Sprintf(levelFmt, STACK)
	msg := fmt.Sprintf(format, v...)
	l.logger.Print(lvl + msg)
}
