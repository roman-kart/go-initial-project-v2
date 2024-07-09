package tests_test

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/roman-kart/go-initial-project/v2/components/tools"
	"github.com/roman-kart/go-initial-project/v2/components/utils"
)

var errTest = errors.New("test error")

func TestPanicOnErrorMustPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	tools.PanicOnError(errTest)
}

func TestPanicOnErrorMustNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()
	tools.PanicOnError(nil)
}

func TestGetRootPath(t *testing.T) {
	p := tools.GetRootPath()
	require.DirExistsf(t, p, "The path does not exist")
	t.Logf("Path: %s", p)
}

func TestExecuteCommandWithOutput(t *testing.T) {
	l := utils.GetZapLogger()
	cmd := exec.Command("echo", "test")
	out, err := tools.ExecuteCommandWithOutput(cmd, l)
	require.NoError(t, err, "Error should be nil")
	assert.Equal(t, "test\n", out, "Output should be 'test'")
}
