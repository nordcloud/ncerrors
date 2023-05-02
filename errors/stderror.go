package errors

import "errors"

// A set of convenience wrappers for standard library 'errors' functions.
var (
	Unwrap = errors.Unwrap
	Is     = errors.Is
	As     = errors.As
)
