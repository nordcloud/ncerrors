package ncerrors

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
				"github.com/nordcloud/ncerrors/errors/error.go(New):76",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(innerFunc):10",
			},
		},
		{
			func() error { return outerFunc() },
			[]string{
				"github.com/nordcloud/ncerrors/errors/error.go(New):76",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(innerFunc):10",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(outerFunc):14",
			},
		},
		{
			func() error { return testStruct{outerFunc}.method() },
			[]string{
				"github.com/nordcloud/ncerrors/errors/error.go(New):76",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(innerFunc):10",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(outerFunc):14",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(testStruct.method):22",
			},
		},
		{
			func() error { return testStruct{innerFunc}.nested() },
			[]string{
				"github.com/nordcloud/ncerrors/errors/error.go(New):76",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(innerFunc):10",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(testStruct.nested.func1):27",
				"github.com/nordcloud/ncerrors/errors/stack_trace_test.go(testStruct.nested):29",
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
	assert.Equal(t, "github.com/nordcloud/ncerrors/errors/stack_trace_test.go(TestGetSingleTrace):78", s[0])
}

func TestGetCallStackTrace(t *testing.T) {
	s := GetTrace()
	assert.Equal(t, 2, len(s))
	assert.Equal(t, "github.com/nordcloud/ncerrors/errors/stack_trace_test.go(TestGetCallStackTrace):85", s[0])
}
