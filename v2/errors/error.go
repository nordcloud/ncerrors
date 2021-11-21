package errors

// NCError is the error type that can store additional message, fields and stack trace.
// It should be used as a standard error.
// If more information about error is required, call GetInfo(err).
type NCError struct {
	err error

	message string
	fields  Fields

	stackTrace stackTrace
	funcName   string
}

// Error returns error message based on all wrapped errors.
func (e NCError) Error() string {
	if e.err != nil {
		return e.message + ": " + e.err.Error()
	}

	return e.message
}

// Info returns info about error.
// It should not be used directly. Use GetInfo(err) instead.
func (e NCError) Info() Info {
	return Info{
		Message:    e.message,
		Fields:     e.fields,
		StackTrace: e.stackTrace.stringStack(),
		FuncName:   e.funcName,
	}
}

// Unwrap unwraps the error.
// It implements the errors.Unwrap interface
func (e NCError) Unwrap() error {
	return e.err
}

// StackTrace returns stack trace of an error.
func (e NCError) StackTrace() stackTrace {
	return e.stackTrace
}

// New creates a new instance of NCError.
func New(message string, fields Fields) NCError {
	var funcName string

	stackTrace := newStackTrace(4)
	if len(stackTrace.Frames) != 0 {
		funcName = stackTrace.Frames[0].functionName.String()
	}

	return NCError{
		message:    message,
		fields:     fields,
		stackTrace: stackTrace,
		funcName:   funcName,
	}
}

// NewWithErr always creates a new instance of NCError, even if err == nil.
// It is required for a proper wrapping with custom error, so that we always get a non-null instance of NCError thus
// avoiding any potential nil pointer dereferences.
// It is similar to both New and Wrap:
//   when err == nil, then it behaves the same as New
//   when err != nil, then it behaves the same as Wrap
func NewWithErr(err error, message string, fields Fields) NCError {
	var funcName string

	stTrace := newStackTrace(4)
	if len(stTrace.Frames) != 0 {
		funcName = stTrace.Frames[0].functionName.String()
	}

	if errorHasStackTrace(err) {
		stTrace = stackTrace{}
	}

	return NCError{
		message:    message,
		fields:     fields,
		err:        err,
		stackTrace: stTrace,
		funcName:   funcName,
	}
}

// Wrap lets you to wrap an error with specyfing message and fields.
// If the err is nil, then the returned error is nil as well.
func Wrap(err error, message string, fields Fields) error {
	if err == nil {
		return nil
	}

	var funcName string

	stTrace := newStackTrace(4)
	if len(stTrace.Frames) != 0 {
		funcName = stTrace.Frames[0].functionName.String()
	}

	if errorHasStackTrace(err) {
		stTrace = stackTrace{}
	}

	return NCError{
		message:    message,
		fields:     fields,
		err:        err,
		stackTrace: stTrace,
		funcName:   funcName,
	}
}

// W lets you to wrap an error without specyfing message or fields.
// The message will be taken from the function name of the function where W is called.
// The fields will remain empty.
func W(err error) error {
	if err == nil {
		return nil
	}

	var message, funcName string

	stTrace := newStackTrace(4)
	if len(stTrace.Frames) != 0 {
		message = stTrace.Frames[0].functionName.name
		funcName = stTrace.Frames[0].functionName.String()
	}

	if errorHasStackTrace(err) {
		stTrace = stackTrace{}
	}

	return NCError{
		message:    message,
		fields:     nil,
		err:        err,
		stackTrace: stTrace,
		funcName:   funcName,
	}
}

func errorHasStackTrace(err error) bool {
	var st stackTracer
	return As(err, &st)
}

type stackTracer interface {
	StackTrace() stackTrace
}
