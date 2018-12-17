package logging

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

type loggingConfig struct {
	jsonLog bool
}

var logConf = loggingConfig{
	jsonLog: false,
}

type Logger struct {
	prefix string
	logger *zap.Logger
}

func NewDevLog(prefix string) *zap.Logger {
	return New(prefix, false)
}

func New(prefix string, isProduction bool) *zap.Logger {
	var err error = nil
	var logger *zap.Logger = nil

	if isProduction {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment(
			zap.AddStacktrace(zap.ErrorLevel),
		)
	}

	if err != nil {
		fmt.Println("can't initialize zap logger: %v", err)
		os.Exit(1)
	}
	defer logger.Sync()

	return logger
}
