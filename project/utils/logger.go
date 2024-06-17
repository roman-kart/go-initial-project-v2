package utils

import (
	"fmt"

	"go.uber.org/zap"

	cfg "github.com/roman-kart/go-initial-project/project/config"
	"github.com/roman-kart/go-initial-project/project/tools"
)

// GetZapLogger returns a [zap.Logger].
func GetZapLogger() *zap.Logger {
	l, err := zap.NewDevelopment()
	tools.PanicOnError(err)

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
			Initial:    config.Logger.Sampling.Initial,
			Thereafter: config.Logger.Sampling.Thereafter,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()

	//nolint:forbidigo
	return &Logger{Config: config, Logger: logger},
		func() { err := logger.Sync(); fmt.Println("Logger sync error:", err) },
		ew(err)
}
