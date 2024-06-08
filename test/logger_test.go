package test_test

import (
	"testing"

	"github.com/roman-kart/go-initial-project/gip"
)

func TestGetLoggerMustNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()

	logger := gip.GetLogger()
	logger.Info("This is a test")
}
