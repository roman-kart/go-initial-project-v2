package tools

import (
	"errors"
	"fmt"
)

// WrapMethodError wrap an error with method name.
func WrapMethodError(err error, method string) error {
	if err == nil {
		return nil
	}

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

// ErrorWrapperCreator creates functions for wrapping errors.
type ErrorWrapperCreator struct {
	Prefix string
}

// NewErrorWrapperCreator creates a new ErrorWrapperCreator.
func NewErrorWrapperCreator() ErrorWrapperCreator {
	return ErrorWrapperCreator{
		Prefix: "",
	}
}

// AppendToPrefix append content to the end of prefix.
// Between content and current prefix added separator.
func (w ErrorWrapperCreator) AppendToPrefix(content string) ErrorWrapperCreator {
	if w.Prefix != "" {
		w.Prefix = w.concatWithPrefix(content)
	} else {
		w.Prefix = content
	}

	return w
}

func (w ErrorWrapperCreator) concatWithPrefix(content string) string {
	return fmt.Sprintf("%s.%s", w.Prefix, content)
}

// GetMethodWrapper returns a function for wrapping errors.
func (w ErrorWrapperCreator) GetMethodWrapper(methodName string) func(err error) error {
	return GetErrorWrapper(w.concatWithPrefix(methodName))
}

// GetErrorWrapper returns a function for wrapping errors.
func GetErrorWrapper(prefix string) func(err error) error {
	return func(err error) error {
		return WrapMethodError(err, prefix)
	}
}
