package logger

import (
	"go.uber.org/zap"
)

var Log zap.Logger

func Load() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	Log = *logger
	return nil
}
