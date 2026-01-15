package utils

import (
	"log"
	"os"
)

var (
    ErrorLogger *log.Logger
    InfoLogger  *log.Logger
    DebugLogger *log.Logger
)

func init() {
    ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
    InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
    DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}