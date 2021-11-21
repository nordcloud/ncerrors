package errors

type NCError struct {
	err error

	message string
	fields  Fields

	stackTrace StackTrace
	funcName   string
}

func (e NCError) Error() string {
	if e.err != nil {
		return e.message + ": " + e.err.Error()
	}

	return e.message
}

func (e NCError) Info() Info {
	return Info{
		Message:    e.message,
		Fields:     e.fields,
		StackTrace: e.stackTrace.StringStack(),
		FuncName:   e.funcName,
	}
}

func New(message string, fields Fields) NCError {
	var funcName string

	stackTrace := newStackTrace(4)
	if len(stackTrace.Frames) != 0 {
		funcName = stackTrace.Frames[0].FunctionName.String()
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

	stackTrace := newStackTrace(4)
	if len(stackTrace.Frames) != 0 {
		funcName = stackTrace.Frames[0].FunctionName.String()
	}

	if errorHasStackTrace(err) {
		stackTrace = StackTrace{}
	}

	return NCError{
		message:    message,
		fields:     fields,
		err:        err,
		stackTrace: stackTrace,
		funcName:   funcName,
	}
}

func Wrap(err error, message string, fields Fields) error {
	if err == nil {
		return nil
	}

	var funcName string

	stackTrace := newStackTrace(4)
	if len(stackTrace.Frames) != 0 {
		funcName = stackTrace.Frames[0].FunctionName.String()
	}

	if errorHasStackTrace(err) {
		stackTrace = StackTrace{}
	}

	return NCError{
		message:    message,
		fields:     fields,
		err:        err,
		stackTrace: stackTrace,
		funcName:   funcName,
	}
}

func W(err error) error {
	if err == nil {
		return nil
	}

	var message, funcName string

	stackTrace := newStackTrace(4)
	if len(stackTrace.Frames) != 0 {
		message = stackTrace.Frames[0].FunctionName.Name
		funcName = stackTrace.Frames[0].FunctionName.String()
	}

	if errorHasStackTrace(err) {
		stackTrace = StackTrace{}
	}

	return NCError{
		message:    message,
		fields:     nil,
		err:        err,
		stackTrace: stackTrace,
		funcName:   funcName,
	}
}

func errorHasStackTrace(err error) bool {
	var st StackTracer
	return As(err, &st)
}

func (e NCError) Unwrap() error {
	return e.err
}

func (e NCError) StackTrace() StackTrace {
	return e.stackTrace
}

type StackTracer interface {
	StackTrace() StackTrace
}

type Infoer interface {
	Info() Info
}

type Unwrapper interface {
	Unwrap() error
}
