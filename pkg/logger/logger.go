package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARNING
	ERROR
	FATAL
)

const DefaultCallerDepth = 2

var logLevels = map[LogLevel]string{
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

func msg(level LogLevel, message string, args ...any) {
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

func Info(message string, args ...any) {
	msg(INFO, message, args...)
}

func Warning(message string, args ...any) {
	msg(WARNING, message, args...)
}

func Error(message string, args ...any) {
	msg(ERROR, message, args...)
}

func Fatal(message string, args ...any) {
	msg(FATAL, message, args...)

	os.Exit(1)
}
