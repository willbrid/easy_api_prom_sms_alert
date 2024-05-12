package logging

import (
	"log"
	"os"
)

type LogLevel int

const (
	Info LogLevel = iota
	Warning
	Error
)

var logPrefix = map[LogLevel]string{
	Info:    "INFO: ",
	Warning: "WARNING: ",
	Error:   "ERROR: ",
}

func Log(level LogLevel, template string, vals ...interface{}) {
	prefix, ok := logPrefix[level]
	if !ok {
		prefix = "UNKNOWN: "
	}
	log.SetOutput(os.Stdout)
	log.SetPrefix(prefix)
	log.Printf(template, vals...)
}
