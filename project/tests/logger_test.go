package tests_test

import (
	"testing"

	"github.com/roman-kart/go-initial-project/project/utils"
)

func TestGetZapLoggerMustNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()

	logger := utils.GetZapLogger()
	logger.Info("This is a test")
}
