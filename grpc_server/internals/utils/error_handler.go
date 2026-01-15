package utils

import (
	"fmt"
)

func ErrorHandler(err error, message string) error {
	if err != nil {
		ErrorLogger.Printf("%v\n", err)
	}
	ErrorLogger.Println(message)
	return fmt.Errorf("%s", message)
}
