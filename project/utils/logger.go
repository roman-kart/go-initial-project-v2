package utils

import (
	cfg "github.com/roman-kart/go-initial-project/project/config"
	"go.uber.org/zap"
)

// GetZapLogger returns a [zap.Logger].
func GetZapLogger() *zap.Logger {
	l, err := zap.NewDevelopment()
	PanicOnError(err)

	return l
}

// Logger is a Logger component of the application.
type Logger struct {
	Config *cfg.Config
	Logger *zap.Logger
}

// NewLogger returns a new Logger component.
func NewLogger(config *cfg.Config) (*Logger, func(), error) {
	logLevel := zap.InfoLevel
	switch config.Logger.Level {
	case "debug":
		logLevel = zap.DebugLevel
	}

	logger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	return &Logger{Config: config, Logger: logger}, func() { logger.Sync() }, err
}
