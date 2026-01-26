package utils

import (
	"fmt"

	"go.uber.org/zap"
)

func ErrorHandler(err error, message string) error {
	if err != nil {
		Logger.Error(message, zap.Error(err))
	} else {
		Logger.Error(message)
	}
	return fmt.Errorf("%s", message)
}
