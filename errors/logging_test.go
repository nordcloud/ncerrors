package errors

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLogger_StadardError(t *testing.T) {
	err := errors.New("Error")
	logEntry := getLogger(err, (*NCError).GetContext)

	assert.NotNil(t, logEntry)
	assert.Equal(t, logrus.Fields{"error": err.Error()}, logEntry.Data)
}

func TestGetLogger_NCError(t *testing.T) {
	errorMessage := "error message"
	errorLevel := "level"

	err := WithContext(errors.New(errorMessage), errorLevel, Fields{"field1": "val1"})
	logEntry := getLogger(err, (*NCError).GetContext)

	assert.NotNil(t, logEntry)
	assert.Contains(t, logEntry.Data, "error_context")

	errCtx := logEntry.Data["error_context"].(Fields)
	assert.Contains(t, errCtx, "causes")

	assert.Equal(t, []Cause{Cause{
		Message:  errorLevel,
		Fields:   Fields{"field1": "val1"},
		FuncName: "TestGetLogger_NCError",
		FileName: "github.com/nordcloud/ncerrors/errors/logging_test.go",
		Line:     25,
		Severity: ERROR,
	}, // this value must be updated according to the line number when the error has actually occured
		Cause{Message: errorMessage}}, errCtx["causes"])
}

func TestGetLogField(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want logrus.Fields
	}{
		{
			name: "No error",
			args: args{err: nil},
			want: logrus.Fields{
				"error": nil,
			},
		},
		{
			name: "Standard error",
			args: args{err: errors.New("dummy")},
			want: logrus.Fields{
				"error": "dummy",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLogFields(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLogFieldAWSError(t *testing.T) {
	err := awserr.New("code", "message", errors.New("org error"))
	logFields := GetLogFields(err)
	assert.True(t, reflect.DeepEqual(logFields, logrus.Fields{
		errorKey:           "org error",
		awsErrorCodeKey:    "code",
		awsErrorMessageKey: "message",
	}))
}

func TestGetLogFieldsNCErrorAWSErrorWrapped(t *testing.T) {
	awsErr := awserr.New("code", "message", errors.New("org error"))
	ncErr := WithContext(awsErr, "context 1", nil)
	logFields := GetLogFields(ncErr)
	assert.Equal(t, logFields[awsErrorCodeKey], "code")
	assert.Equal(t, logFields[awsErrorMessageKey], "message")
}

func TestGetLogField_NCError(t *testing.T) {
	errorMessage := "dummy"
	errorLevel := "level"

	err := WithContext(errors.New(errorMessage), errorLevel, Fields{"field1": "val1"})
	logFields := GetLogFields(err)

	assert.NotNil(t, logFields)
	assert.Contains(t, logFields, "error_context")

	errCtx := logFields["error_context"].(Fields)
	assert.Contains(t, errCtx, "causes")

	assert.Equal(t, []Cause{Cause{
		Message:  errorLevel,
		Fields:   Fields{"field1": "val1"},
		FuncName: "TestGetLogField_NCError",
		FileName: "github.com/nordcloud/ncerrors/errors/logging_test.go",
		Line:     100,
		Severity: ERROR,
	}, // this value must be updated according to the line number when the error has actually occured
		Cause{Message: errorMessage}}, errCtx["causes"])
}

func TestGetMergedLogField_NCError(t *testing.T) {
	errorMessage := "dummy"
	errorLevel := "level"

	err1 := WithContext(errors.New(errorMessage), errorLevel, Fields{
		"field1": "val1",
		"field2": "val2",
	})
	err2 := WithContext(err1, errorMessage, Fields{"field2": "override"})
	logFields := GetMergedLogFields(err2)

	assert.NotNil(t, logFields)
	assert.Contains(t, logFields, "error_context")

	errCtx := logFields["error_context"].(Fields)
	assert.Contains(t, errCtx, "fields")

	assert.Equal(t, Fields{
		"field1": "val1",
		"field2": "override",
	}, errCtx["fields"])
}
