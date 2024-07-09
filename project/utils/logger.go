package utils

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/roman-kart/go-initial-project/v2/project/tools"
)

// GetZapLogger returns a [zap.Logger].
func GetZapLogger() *zap.Logger {
	l, err := zap.NewDevelopment()
	tools.PanicOnError(err)

	return l
}

type LoggerConfig struct {
	Level    string
	Sampling struct {
		Initial    int
		Thereafter int
	}
}

// NewLogger returns a new Logger component.
func NewLogger(config *LoggerConfig) (*zap.Logger, func(), error) {
	logLevel := zap.InfoLevel

	switch config.Level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "panic":
		logLevel = zap.PanicLevel
	}

	ew := tools.GetErrorWrapper("NewUserManager")

	logger, err := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    config.Sampling.Initial,
			Thereafter: config.Sampling.Thereafter,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()

	//nolint:forbidigo
	return logger,
		func() { err := logger.Sync(); fmt.Println("Logger sync error:", err) },
		ew(err)
}
