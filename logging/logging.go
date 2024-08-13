package logging

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
	Info LogLevel = iota
	Warning
	Error
)

const DefaultCallerDepth = 1

var logLevels = map[LogLevel]string{
	Info:    "INFO",
	Warning: "WARNING",
	Error:   "ERROR",
}

func Log(level LogLevel, template string, vals ...interface{}) {
	logLevel, ok := logLevels[level]
	if !ok {
		logLevel = "UNKNOWN"
	}

	datetime := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(template, vals...)

	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	filename := ""
	if ok {
		filename = filepath.Base(file)
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	log.Printf("[%s][%s][%s:%d] : %s", datetime, logLevel, filename, line, message)
}
