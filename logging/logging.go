package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	Info LogLevel = iota
	Warning
	Error
)

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

	pc, _, _, ok := runtime.Caller(1)
	functionName := "unknown"
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		parts := strings.Split(functionName, "/")
		functionName = parts[len(parts)-1]
	}

	message := fmt.Sprintf(template, vals...)

	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	log.Printf("[%s][%s][%s] : %s", datetime, logLevel, functionName, message)
}
