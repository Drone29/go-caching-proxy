package logger

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	err   *log.Logger
	debug *log.Logger
}

func (logger *Logger) Infof(format string, v ...any) {
	logger.info.Printf(format, v...)
}

func (logger *Logger) Errorf(format string, v ...any) {
	logger.err.Printf(format, v...)
}

func (logger *Logger) Debugf(format string, v ...any) {
	logger.debug.Printf(format, v...)
}

// create a new logger
func New(debug bool) *Logger {
	logger := &Logger{
		info: log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lmsgprefix),
		err:  log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lmsgprefix),
	}
	if debug {
		logger.debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lmsgprefix)
	}
	return logger
}
