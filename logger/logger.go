package logger

import (
	"log"
	"os"
)

var (
	Inf *log.Logger
	Err *log.Logger
	Dbg *log.Logger
)

func init() {
	Inf = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	Dbg = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime)
	Err = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime)
}
