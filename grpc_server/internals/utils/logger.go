package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() (*zap.Logger, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	cfgPath := filepath.Join(dir, "log_config.json")

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, ErrorHandler(err, "Failed to initialize logger")
	}

	var cfg zap.Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, ErrorHandler(err, "Failed to unmarshal log config")
	}

	logger := zap.Must(cfg.Build())

	return logger, nil
}
