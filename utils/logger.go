package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitializeLogger(verbose bool) {
	config := zap.Config{
		Development: false,
		DisableCaller: true,
		DisableStacktrace: false,
		Encoding: "console",
		EncoderConfig: zap.NewDevelopmentEncoderConfig(),
		Level: zap.NewAtomicLevelAt(zap.InfoLevel),
		Sampling: nil,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	
	if verbose {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	Logger, _ = config.Build()

	// Flush buffer before exiting the application.
	defer Logger.Sync()
}
