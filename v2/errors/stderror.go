package errors

import "errors"

// A set of convinience wrappers for standard library 'errors' functions.
var (
	Unwrap = errors.Unwrap
	Is     = errors.Is
	As     = errors.As
)
