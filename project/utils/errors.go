package utils

import (
	"errors"
	"fmt"
)

// WrapMethodError wrap an error with method name.
func WrapMethodError(err error, method string) error {
	return fmt.Errorf("%s: %w", method, err)
}

// ErrHTTPWrongStatus is an error type for wrong HTTP status.
var ErrHTTPWrongStatus = errors.New("wrong status")

// NewErrHTTPWrongStatus is a function for creating [ErrHTTPWrongStatus].
//
// Accepts two parameters:
//   - expected - expected HTTP status
//   - got - got HTTP status
func NewErrHTTPWrongStatus(expected int, got int) error {
	return fmt.Errorf("%w: expected %d got %d", ErrHTTPWrongStatus, expected, got)
}
