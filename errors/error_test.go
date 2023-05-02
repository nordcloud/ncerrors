// Copyright 2023 Nordcloud Oy or its affiliates. All Rights Reserved.

package errors

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorsListToError(t *testing.T) {
	err := ListToError([]error{})
	assert.Nil(t, err)

	err = ListToError([]error{errors.New("e1"), errors.New("e2")})
	assert.NotNil(t, err)
	assert.Equal(t, "[e1, e2]", err.Error())

}

func f3() error {
	return f2()
}

func f2() error {
	return f1()
}

func f1() error {
	return errors.WithStack(errors.New("Example f1"))
}

func TestErrorWrap(t *testing.T) {
	errorMessage := "error message"
	err := errors.New(errorMessage)

	messageLevel1 := "level1"
	fieldsLevel1 := Fields{"field1": "val1", "field2": 2}
	level1 := WithContext(err, messageLevel1, fieldsLevel1)

	e, _ := level1.(NCError)

	assert.Equal(t, fmt.Sprintf("level1: %s", errorMessage), level1.Error())
	assert.Len(t, e.Stack, 3)
	level1Causes := []Cause{Cause{
		Message:  messageLevel1,
		FuncName: "TestErrorWrap",
		Line:     41,
		FileName: "github.com/nordcloud/ncerrors/errors/error_test.go",
		Fields:   fieldsLevel1,
		Severity: ERROR,
	},
		Cause{Message: errorMessage}}
	errorFields := e.GetContext()
	assert.Equal(t, level1Causes, errorFields["causes"])

	// Next wrapping level. Error from the level1.
	fieldsLevel2 := Fields{"field3": "val2"}
	messageLevel2 := "level2"
	level2 := WithContext(level1, messageLevel2, fieldsLevel2)
	e, _ = level2.(NCError)

	assert.Equal(t, fmt.Sprintf("level2: level1: %s", errorMessage), level2.Error())
	assert.Equal(t, 3, len(e.Stack))
	level2Causes := []Cause{
		Cause{
			Message:  messageLevel2,
			FuncName: "TestErrorWrap",
			Line:     62,
			FileName: "github.com/nordcloud/ncerrors/errors/error_test.go",
			Fields:   fieldsLevel2,
			Severity: ERROR,
		},
		Cause{
			Message:  messageLevel1,
			FuncName: "TestErrorWrap",
			Line:     41,
			FileName: "github.com/nordcloud/ncerrors/errors/error_test.go",
			Fields:   fieldsLevel1,
			Severity: ERROR,
		},
		Cause{Message: errorMessage}}
	errorFields = e.GetContext()
	assert.Equal(t, level2Causes, errorFields["causes"])
}

func TestErrorWrap_preservedStack(t *testing.T) {
	errorMessage := "error message"
	err := errors.New(errorMessage)
	level1 := WithContext(err, "level1", Fields{"field1": "val1", "field2": 2})
	e, _ := level1.(NCError)

	// Stack should be preserved from even if the error is wrapped in the higher levels.
	level2 := WithContext(level1, "level2", Fields{"field3": "val2"})
	e, _ = level2.(NCError)
	assert.Equal(t, fmt.Sprintf("level2: level1: %s", errorMessage), level2.Error())
	assert.Equal(t, 3, len(e.Stack))
}

func TestErrorWrap_preservedRootError(t *testing.T) {
	errorMessage := "error message"
	err := errors.New(errorMessage)

	level1 := WithContext(err, "level1", nil)
	e, _ := level1.(NCError)

	level2 := WithContext(level1, "level2", nil)
	e, _ = level2.(NCError)

	//Check if the RootError is preserved.
	assert.Equal(t, err, e.RootError)
}

func TestUpdateCauses(t *testing.T) {
	initialErrorMessage := "error"
	level1 := WithContext(errors.New(initialErrorMessage), "level1", nil)
	level2 := WithContext(level1, "level2", Fields{"field1": "val1"})
	e, _ := level2.(NCError)
	expectedCauses := []Cause{
		Cause{
			Message:  "level2",
			Fields:   Fields{"field1": "val1"},
			FuncName: "TestUpdateCauses",
			FileName: "github.com/nordcloud/ncerrors/errors/error_test.go",
			Line:     119,
			Severity: ERROR,
		},
		Cause{
			Message:  "level1",
			Fields:   Fields(nil),
			FuncName: "TestUpdateCauses",
			FileName: "github.com/nordcloud/ncerrors/errors/error_test.go",
			Line:     118,
			Severity: ERROR,
		},
		Cause{Message: initialErrorMessage}}
	assert.Equal(t, expectedCauses, e.Causes)
}

func TestNewSimpleError(t *testing.T) {
	errorMessage := "simple error"
	err := New(errorMessage, nil)
	assert.Equal(t, errorMessage, err.Error())
	e, _ := err.(NCError)
	causes := e.Causes
	assert.Len(t, causes, 1)
	assert.Equal(t, errorMessage, causes[0].Message)
	assert.Nil(t, causes[0].Fields)
	assert.Equal(t, "github.com/nordcloud/ncerrors/errors/error_test.go", causes[0].FileName)
	assert.Equal(t, "TestNewSimpleError", causes[0].FuncName)
}

func TestNewMethodWrapError(t *testing.T) {
	err := New("simple error", nil)

	level2 := WithContext(err, "level2", nil)
	level3 := WithContext(level2, "level3", nil)

	assert.Equal(t, "level3: level2: simple error", level3.Error())
}

func TestErrorSeverity(t *testing.T) {
	err := NewWithSeverity("error1", nil, WARN)
	level2 := WithContextAndSeverity(err, "level2", ERROR, nil)
	level3 := WithContextAndSeverity(level2, "level3", DEBUG, nil)

	severity := GetErrorSeverity(level3)
	assert.Equal(t, severity, DEBUG) // outermost severity
	if ncE, ok := level3.(NCError); ok {
		assert.Equal(t, 3, len(ncE.Causes))
	}
}

func TestErrorSeverity_Single(t *testing.T) {
	err := NewWithSeverity("error1", nil, WARN)
	severity := GetErrorSeverity(err)
	assert.Equal(t, WARN, severity)
}

func TestGetRootError(t *testing.T) {
	rootError := errors.New("Root error")
	wrappedError := WithContext(rootError, "second error", Fields{"123": "456"})
	ncErrorNew := New("Root error", nil)

	testCases := []error{
		rootError,
		wrappedError,
		WithContext(wrappedError, "third error", nil),
		ncErrorNew,
	}

	for _, testCase := range testCases {
		res := GetRootError(testCase)
		if res != nil {
			assert.Equal(t, rootError.Error(), res.Error())
		} else {
			assert.Nil(t, res)
		}
	}
}

func TestFieldsAdd(t *testing.T) {
	fields := Fields{"key1": "val1", "key2": "val2"}
	extFields := fields.Add("key3", "val3")
	expectedFields := Fields{"key1": "val1", "key2": "val2", "key3": "val3"}
	assert.Equal(t, expectedFields, extFields)

	// Check if the original fields remain intact.
	assert.Equal(t, nil, fields["key3"])
}

func TestFieldsExtend(t *testing.T) {
	fields := Fields{"key1": "val1", "key2": "val2"}
	extFields := fields.Extend(Fields{"key3": "val3", "key4": "val4"})
	expectedFields := Fields{"key1": "val1", "key2": "val2", "key3": "val3", "key4": "val4"}
	assert.Equal(t, expectedFields, extFields)

	// Check if the original fields remain intact.
	assert.Equal(t, nil, fields["key3"])
	assert.Equal(t, nil, fields["key4"])
}

func TestWrap(t *testing.T) {
	errorMessage := "error message"
	err := errors.New(errorMessage)

	messageLevel1 := "level1"
	fieldsLevel1 := Fields{"field1": "val1", "field2": 2}
	level1 := Wrap(err, messageLevel1, fieldsLevel1)

	assert.Equal(t, fmt.Sprintf("level1: %s", errorMessage), level1.Error())
}

func TestWrap_NilError(t *testing.T) {
	messageLevel1 := "level1"
	fieldsLevel1 := Fields{"field1": "val1", "field2": 2}
	err := Wrap(nil, messageLevel1, fieldsLevel1)

	assert.Nil(t, nil, err)
}
