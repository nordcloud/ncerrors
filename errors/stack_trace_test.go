// Copyright 2023 Nordcloud Oy or its affiliates. All Rights Reserved.

package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func innerFunc() error {
	return New("test", nil)
}

func outerFunc() error {
	return innerFunc()
}

type testStruct struct {
	fn func() error
}

func (ts testStruct) method() error {
	return ts.fn()
}

func (ts testStruct) nested() error {
	fn := func() error {
		return ts.fn()
	}
	return fn()
}

func TestSimpleFuncStack(t *testing.T) {
	for _, tc := range []struct {
		fn    func() error
		stack []string
	}{
		{
			func() error { return innerFunc() },
			[]string{
				"github.com/nordcloud/ncerrors/errors/error.go(New):135",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(innerFunc):12",
			},
		},
		{
			func() error { return outerFunc() },
			[]string{
				"github.com/nordcloud/ncerrors/errors/error.go(New):135",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(innerFunc):12",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(outerFunc):16",
			},
		},
		{
			func() error { return testStruct{outerFunc}.method() },
			[]string{
				"github.com/nordcloud/ncerrors/errors/error.go(New):135",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(innerFunc):12",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(outerFunc):16",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(testStruct.method):24",
			},
		},
		{
			func() error { return testStruct{innerFunc}.nested() },
			[]string{
				"github.com/nordcloud/ncerrors/errors/error.go(New):135",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(innerFunc):12",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(testStruct.nested.func1):29",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(testStruct.nested):31",
			},
		},
	} {
		err := tc.fn()
		ncErr := err.(NCError)
		assert.Equal(t, tc.stack, ncErr.Stack[:len(tc.stack)])
	}
}

func TestGetSingleTrace(t *testing.T) {
	s := GetTrace()
	// Returns list of the stack trace
	assert.Len(t, s, 2)
	assert.Equal(t, "github.com/nordcloud/ncerrors/errors/stack_trace_test.go(TestGetSingleTrace):80", s[0])
}

func TestGetCallStackTrace(t *testing.T) {
	s := GetTrace()
	assert.Len(t, s, 2)
	assert.Equal(t, "github.com/nordcloud/ncerrors/errors/stack_trace_test.go(TestGetCallStackTrace):87", s[0])
}
