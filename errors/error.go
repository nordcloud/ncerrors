package errors

import (
	"fmt"
	"strings"
)

type LogSeverity string

const (
	//ERROR Severity - logged with error level
	ERROR LogSeverity = "error"
	//WARN Severity warn - logged with warning level
	WARN LogSeverity = "warning"
	//INFO Severity - logged with info level
	INFO LogSeverity = "info"
	//DEBUG Severity - logged with Debug level
	DEBUG LogSeverity = "debug"
)

// ListToError converts errors list to single error
func ListToError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var errorMessages []string
	for _, err := range errs {
		errorMessages = append(errorMessages, err.Error())
	}
	return fmt.Errorf("[%s]", strings.Join(errorMessages[:], ", "))
}

// Fields keeps context.
type Fields map[string]interface{}

// Add adds key:val to the Fields, returns fresh extended copy. The original Fields remains intact.
func (f Fields) Add(key string, val interface{}) Fields {
	newFields := f.copy()
	newFields[key] = val
	return newFields
}

// Extend extends Fields with the content of extFields, returns fresh extended copy. The original Fields remains intact.
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

// Cause keeps the context information about the error.
type Cause struct {
	Message  string
	Fields   Fields
	FuncName string
	FileName string
	Line     int
	Severity LogSeverity
}

// NCError basic error structure.
type NCError struct {
	Causes []Cause
	// Contains stack trace from the initial place when the error
	// was raised.
	Stack []string
	//The root error at the base level.
	RootError error
}

func (n NCError) Error() string {
	var messages []string
	for _, cause := range n.Causes {
		messages = append(messages, cause.Message)
	}
	return strings.Join(messages, ": ")
}

// New error with context.
func New(message string, fields Fields) error {
	fileName, funcName, lineNumber := GetRuntimeContext()
	newCause := Cause{
		Message:  message,
		Fields:   fields,
		FuncName: funcName,
		FileName: fileName,
		Line:     lineNumber,
		Severity: ERROR}
	return NCError{
		Causes: []Cause{newCause},
		Stack:  GetTrace()}
}

func NewWithSeverity(message string, fields Fields, severity LogSeverity) error {
	fileName, funcName, lineNumber := GetRuntimeContext()
	newCause := Cause{
		Message:  message,
		Fields:   fields,
		FuncName: funcName,
		FileName: fileName,
		Line:     lineNumber,
		Severity: severity,
	}
	return NCError{
		Causes: []Cause{newCause},
		Stack:  GetTrace(),
	}
}

//WithContext set new error wrapped with message and error context.
func WithContext(err error, message string, fields Fields) error {
	// Attach message to the list of causes.
	fileName, funcName, lineNumber := GetRuntimeContext()
	newCause := Cause{
		Message:  message,
		Fields:   fields,
		FuncName: funcName,
		FileName: fileName,
		Line:     lineNumber,
		Severity: ERROR,
	}
	//If we wrap existing NCError at the higher layer. Here we only append causes.
	//and do not touch stack trace and root error.
	if ncError, ok := err.(NCError); ok {
		ncError.Causes = append([]Cause{newCause}, ncError.Causes...)
		return ncError
	}

	return NCError{
		Causes:    []Cause{newCause, Cause{Message: err.Error()}},
		Stack:     GetTrace(),
		RootError: err}
}

//WithContextAndSeverity set new error wrapped with message, severity and error context.
func WithContextAndSeverity(err error, message string, severity LogSeverity, fields Fields) error {
	// Attach message to the list of causes.
	fileName, funcName, lineNumber := GetRuntimeContext()
	newCause := Cause{
		Message:  message,
		Fields:   fields,
		FuncName: funcName,
		FileName: fileName,
		Line:     lineNumber,
		Severity: severity,
	}
	//If we wrap existing NCError at the higher layer. Here we only append causes.
	//and do not touch stack trace and root error.
	if ncError, ok := err.(NCError); ok {
		ncError.Causes = append([]Cause{newCause}, ncError.Causes...)
		return ncError
	}

	return NCError{
		Causes:    []Cause{newCause, Cause{Message: err.Error()}},
		Stack:     GetTrace(),
		RootError: err,
	}
}

// GetContext returns fields from the error (with attached stack and causes fields)
// This will be used for logrus.WithFields method.
func (n *NCError) GetContext() Fields {
	return Fields{
		"stack":  n.Stack,
		"causes": n.Causes}
}

// GetMergedFields returns fields from the error.
// The custom fields are merged from every error's cause's fields (as defined by WithContext invocations).
// If the field with the same name is present in multiple causes the value from the outermost cause is taken.
// This will be used for logrus.WithFields method.
func (n *NCError) GetMergedFields() Fields {
	errFields := Fields{}
	for _, cause := range n.Causes {
		for k, v := range cause.Fields {
			if _, exists := errFields[k]; !exists {
				errFields[k] = v
			}
		}
	}

	return errFields
}

// GetMergedFieldsContext returns error stack and merged fields.
func (n *NCError) GetMergedFieldsContext() Fields {
	return Fields{
		"stack":  n.Stack,
		"fields": n.GetMergedFields(),
	}
}

//GetErrorSeverity returns outermost NCError severity or ERROR level.
func GetErrorSeverity(err error) LogSeverity {
	if ncError, ok := err.(NCError); ok {
		if len(ncError.Causes) > 0 {
			return ncError.Causes[0].Severity
		}
		return ERROR
	}
	return ERROR
}

// GetRootError returns root error.
func GetRootError(err error) error {
	if ncError, ok := err.(NCError); ok && ncError.RootError != nil {
		return ncError.RootError
	}
	return err
}

// Unwrap returns underlying non-NCError error. If no errors were wrapped by NCError then returns nil.
func (n NCError) Unwrap() error {
	return n.RootError
}

// Wrap wraps WithContext and checks for nil error.
func Wrap(err error, message string, fields Fields) error {
	if err == nil {
		return nil
	}

	return WithContext(err, message, fields)
}
