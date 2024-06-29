package tools

import (
	"context"
	"crypto/rand"
	"io"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PanicOnError panics if err is not nil.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// GetRootPath returns the root path of the project.
// Panic on error.
func GetRootPath() string {
	rootPath, err := os.Getwd()
	PanicOnError(err)

	rootPathStr := strings.TrimSpace(rootPath)

	return rootPathStr
}

// GetPathFromRoot returns the path of the given path relative to the project root path.
func GetPathFromRoot(path string) string {
	return filepath.Join(GetRootPath(), path)
}

// ExecuteCommandWithOutput executes the given command and returns the output.
//
// Accepts:
//   - cmd *[exec.Cmd] - command for executing
//   - logger *[zap.Logger] - logger
//
// Returns:
//   - string - output of the command
//   - error - nil or error if any error occurred
func ExecuteCommandWithOutput(cmd *exec.Cmd, logger *zap.Logger) (string, error) {
	logger = logger.Named("ExecuteCommandWithOutput")

	cmdString := cmd.String()
	logger = logger.With(zap.String("command", cmdString))
	output, err := cmd.Output()
	outputStr := string(output)

	if err != nil {
		logger.Error("Error while executing command", zap.Error(err), zap.String("commandOutput", outputStr))
		return "", WrapMethodError(err, "ExecuteCommandWithOutput")
	}

	logger.Info("Command executed", zap.String("commandOutput", outputStr))

	return outputStr, nil
}

// GenerateRandomInt returns a random integer between min and max.
func GenerateRandomInt(min, max int) int {
	return int(GenerateRandomInt64(int64(min), int64(max)))
}

// GenerateRandomInt64 returns a random int64 between min and max.
func GenerateRandomInt64(min, max int64) int64 {
	if min >= max {
		return min
	}

	offset := max - min + 1
	randomBigInt, err := rand.Int(rand.Reader, big.NewInt(offset))
	PanicOnError(err)

	return min + randomBigInt.Int64()
}

// DownloadFile download a file from URL to specific filepath
//
// Accepts:
//   - filepath - path to the file
//   - url - URL of the file
//   - logger *[zap.Logger] - logger
//
// Returns:
//   - error - nil or error if any error occurred
func DownloadFile(filepath string, url string, logger *zap.Logger) error {
	return DownloadFileWithContext(context.Background(), filepath, url, logger)
}

// DownloadFileWithContext download a file from URL to specific filepath with context
//
// Accepts:
//   - ctx context.Context - context
//   - filepath - path to the file
//   - url - URL of the file
//   - logger *[zap.Logger] - logger
//
// Returns:
//   - error - nil or error if any error occurred
func DownloadFileWithContext(ctx context.Context, filepath string, url string, logger *zap.Logger) error {
	ew := GetErrorWrapper("DownloadFileWithContext")
	logger = logger.Named("DownloadFileWithContext")

	logger.Info("Downloading file", zap.String("url", url), zap.String("filepath", filepath))

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return WrapMethodError(err, "DownloadFileWithContext")
	}
	defer out.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return WrapMethodError(err, "DownloadFileWithContext")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return WrapMethodError(err, "DownloadFileWithContext")
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return ew(err)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return ew(err)
	}

	return nil
}

// GenerateUUID generate a UUID string.
func GenerateUUID() string {
	u := uuid.New()
	return u.String()
}

// SortMapKeys sort map's keys.
func SortMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

// FirstNonEmpty returns the first non-empty value.
// If no non-empty value is found, returns the default value.
// nolint: ireturn
func FirstNonEmpty[T comparable](values ...T) T {
	var emptyValue T

	for _, value := range values {
		if value != emptyValue {
			return value
		}
	}

	return emptyValue
}
