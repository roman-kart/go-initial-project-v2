package test_test

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/roman-kart/go-initial-project/gip"
)

var errTest = errors.New("test error")

func TestPanicOnErrorMustPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	gip.PanicOnError(errTest)
}

func TestPanicOnErrorMustNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()
	gip.PanicOnError(nil)
}

func TestGetRootPath(t *testing.T) {
	p, err := gip.GetRootPath()
	require.NoError(t, err, "Error should be nil")
	assert.DirExistsf(t, p, "The path does not exist")
	t.Logf("Path: %s", p)
}

func TestExecuteCommandWithOutput(t *testing.T) {
	l := gip.GetLogger()
	cmd := exec.Command("echo", "test")
	out, err := gip.ExecuteCommandWithOutput(cmd, l)
	require.NoError(t, err, "Error should be nil")
	assert.Equal(t, "test\n", out, "Output should be 'test'")
}
