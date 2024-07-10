package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/roman-kart/go-initial-project/v2/components/tools"
)

// GetZapLogger returns a [zap.Logger].
func GetZapLogger() *zap.Logger {
	l, err := zap.NewDevelopment()
	tools.PanicOnError(err)

	return l
}

type LoggerConsoleConfig struct {
	Level string
}

type LoggerRotationConfig struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Localtime  bool
	Compress   bool
}

type LoggerFileConfig struct {
	Level    string
	Path     string
	Rotation *LoggerRotationConfig
}

type LoggerConfig struct {
	Console *LoggerConsoleConfig
	File    *LoggerFileConfig
}

func NewProductionEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg
}

func NewConsoleZapCore(config *LoggerConsoleConfig) zapcore.Core {
	level := ConvertZapLevel(config.Level)

	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(level),
	)
}

func NewFileZapCore(config *LoggerFileConfig) (zapcore.Core, error) {
	level := ConvertZapLevel(config.Level)
	ew := tools.GetErrorWrapper("NewUserManager")

	dir := filepath.Dir(config.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, ew(err)
	}

	var writer io.Writer

	writer, err := os.OpenFile(config.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, ew(err)
	}

	if config.Rotation != nil {
		writer = &lumberjack.Logger{
			Filename:   config.Path,
			MaxSize:    config.Rotation.MaxSize,
			MaxAge:     config.Rotation.MaxAge,
			MaxBackups: config.Rotation.MaxBackups,
			LocalTime:  config.Rotation.Localtime,
			Compress:   config.Rotation.Compress,
		}
	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(NewProductionEncoderConfig()),
		zapcore.AddSync(writer),
		zap.NewAtomicLevelAt(level),
	), nil
}

// NewLogger returns a new Logger component.
func NewLogger(config *LoggerConfig) (*zap.Logger, func(), error) {
	ew := tools.GetErrorWrapper("NewUserManager")
	cores := []zapcore.Core{}

	if config.Console != nil {
		cores = append(cores, NewConsoleZapCore(config.Console))
	}

	if config.File != nil {
		core, err := NewFileZapCore(config.File)
		if err != nil {
			return nil, nil, ew(err)
		}

		cores = append(cores, core)
	}

	logger := zap.New(zapcore.NewTee(cores...))

	//nolint:forbidigo
	return logger,
		func() { err := logger.Sync(); fmt.Println("Logger sync error:", err) },
		nil
}

// ConvertZapLevel converts a string level to a zapcore.Level.
// e.g. "debug" -> zapcore.DebugLevel.
// If level is invalid, it returns zapcore.InfoLevel.
func ConvertZapLevel(level string) zapcore.Level {
	logLevel := zapcore.InfoLevel

	switch level {
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

	return logLevel
}
