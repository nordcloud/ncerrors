package errors

import (
	"errors"
	"fmt"
)

type Fields map[string]interface{}

func (f Fields) Add(key string, val interface{}) Fields {
	newFields := f.copy()
	newFields[key] = val

	return newFields
}

func (f Fields) Extend(extFields Fields) Fields {
	newFields := f.copy()
	for k, v := range extFields {
		newFields[k] = v
	}

	return newFields
}

func (f Fields) copy() Fields {
	newFields := make(map[string]interface{}, len(f))
	for k, v := range f {
		newFields[k] = v
	}

	return newFields
}

type Infoer interface {
	Info() Info
}

type Unwrapper interface {
	Unwrap() error
}

// A collection of all things we know about the error.
type Info struct {
	Message    string   `json:"message,omitempty"`
	Fields     Fields   `json:"fields,omitempty"`
	StackTrace []string `json:"stackTrace,omitempty"`
	FuncName   string   `json:"funcName,omitempty"`
}

func GetInfo(err error) []Info {
	if err == nil {
		return nil
	}

	var infos []Info

	for err != nil {
		// This one should be type asserted instead of using errors.As,
		// because in case err does not implement Info but instead implements Unwrap
		// we can get an info for the unwrapped error instead
		if infoer, ok := err.(Infoer); ok {
			infos = append(infos, infoer.Info())
		} else {
			infos = append(infos, Info{
				Message: err.Error(),
			})
		}

		var unwrapper Unwrapper
		// it is safe to use errors.As here, since if the err implements Unwrap, we will get it directly
		// and if it dont implement Unwrap it wont get unwrapped anyway
		if errors.As(err, &unwrapper) {
			err = unwrapper.Unwrap()
		} else {
			err = nil
		}
	}

	return infos
}

type NCError struct {
	err error

	message string
	fields  Fields

	stackTrace StackTrace
	funcName   string
}

func (e NCError) GetFields() Fields {
	return e.fields
}

func (e NCError) Format(s fmt.State, verb rune) {
	//	_, _ = io.WriteString(s, e.message)
	//	_, _ = io.WriteString(s, "\n")
	//	if e.st != nil {
	//		e.st.Format(s, verb)
	//	}
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

func errorHasStackTrace(err error) bool {
	var ncErr NCError
	return As(err, &ncErr)
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

func (e NCError) Unwrap() error {
	return e.err
}

// A set of convinience wrappers for standard library 'errors' functions.
var (
	Unwrap = errors.Unwrap
	Is     = errors.Is
	As     = errors.As
)
