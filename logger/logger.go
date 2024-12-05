package logger

import (
	"log"
	"os"
)

func New(level string) *log.Logger {
	var prefix string
	var out = os.Stdout

	switch level {
	case "INFO", "info":
		prefix = "[INFO] "
	case "DEBUG", "debug":
		prefix = "[DEBUG] "
	case "ERROR", "error":
		prefix = "[ERROR] "
		out = os.Stderr
	default:
		prefix = level
	}
	return log.New(out, prefix, log.Ldate|log.Ltime|log.Lmsgprefix)
}
