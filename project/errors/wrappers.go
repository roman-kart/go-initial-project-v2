package errors

import "fmt"

// WrapMethodError wrap an error with method name.
func WrapMethodError(err error, method string) error {
	return fmt.Errorf("%s: %w", method, err)
}
