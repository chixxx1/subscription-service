package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func NewLogger(env string) (*zap.Logger, error) {
	var config zap.Config

	if env == "dev" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()
	if err != nil {
		err = fmt.Errorf("panic to initialized zap logger: %w", err)
		panic(err)
	}

	return logger, nil
}
