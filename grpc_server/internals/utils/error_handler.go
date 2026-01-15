package utils

import (
	"fmt"
)

func ErrorHandler(err error, message string) error {
	ErrorLogger.Println(message)
	return fmt.Errorf("%s", message)
}
