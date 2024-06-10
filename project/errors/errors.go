package errors

import (
	"errors"
	"fmt"
)

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
