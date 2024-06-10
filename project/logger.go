package project

import "go.uber.org/zap"

// GetLogger returns a [zap.Logger].
func GetLogger() *zap.Logger {
	l, err := zap.NewDevelopment()
	PanicOnError(err)

	return l
}
