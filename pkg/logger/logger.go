package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type ILogger interface {
	Info(message string, args ...any)
	Warning(message string, args ...any)
	Error(message string, args ...any)
	Fatal(message string, args ...any)
}

type Logger struct {
	logger *log.Logger
}

type LogLevel int

const (
	Info LogLevel = iota
	Warning
	Error
	Fatal
)

const DefaultCallerDepth = 2

var logLevels = map[LogLevel]string{
	Info:    "INFO",
	Warning: "WARNING",
	Error:   "ERROR",
	Fatal:   "FATAL",
}

func NewLogger() *Logger {
	logger := log.New(os.Stdout, "", 0)

	return &Logger{logger: logger}
}

func (l *Logger) msg(level LogLevel, message string, args ...any) {
	logLevel, ok := logLevels[level]
	if !ok {
		logLevel = "UNKNOWN"
	}

	datetime := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	filename := ""
	if ok {
		filename = filepath.Base(file)
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	log.Printf("[%s][%s][%s:%d] : %s", datetime, logLevel, filename, line, fmt.Sprintf(message, args...))
}

func (l *Logger) Info(message string, args ...any) {
	l.msg(Info, message, args...)
}

func (l *Logger) Warning(message string, args ...any) {
	l.msg(Warning, message, args...)
}

func (l *Logger) Error(message string, args ...any) {
	l.msg(Error, message, args...)
}

func (l *Logger) Fatal(message string, args ...any) {
	l.msg(Fatal, message, args...)

	os.Exit(1)
}
