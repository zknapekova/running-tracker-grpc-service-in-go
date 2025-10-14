package utils

import (
	"fmt"
	"log"
	"os"
)

func ErrorHandler(err error, message string) error {
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger.Println(message)
	return fmt.Errorf(message)
}