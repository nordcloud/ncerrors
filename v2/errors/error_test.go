package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var rootErr = fmt.Errorf("rootErr")

func dummyFunc() error {
	return W(rootErr)
}

func Test_Wrap(t *testing.T) {
	t.Run("Wrapped error returns concatenation of messages", func(t *testing.T) {
		wrappedErr := Wrap(rootErr, "message", nil)

		assert.Equal(t, "message: rootErr", wrappedErr.Error())
	})

	t.Run("Wrapped error is unwrappable", func(t *testing.T) {
		wrappedErr := Wrap(rootErr, "message", nil)

		assert.Equal(t, rootErr, Unwrap(wrappedErr))
	})
}

func Test_GetInfo(t *testing.T) {
	t.Run("Sentinel errors returns only message", func(t *testing.T) {
		infos := GetInfo(rootErr)

		require.Len(t, infos, 1)
		assert.Equal(t, "rootErr", infos[0].Message)
		assert.Empty(t, infos[0].Fields)
		assert.Empty(t, infos[0].StackTrace)
		assert.Empty(t, infos[0].FuncName)
	})

	t.Run("New NCError errors returns message", func(t *testing.T) {
		rootNCErr := New("rootErr", nil)
		infos := GetInfo(rootNCErr)

		require.Len(t, infos, 1)
		assert.Equal(t, "rootErr", infos[0].Message)
	})

	t.Run("New NCError errors returns fields", func(t *testing.T) {
		rootErr := New("rootErr", Fields{
			"key1": "val1",
			"key2": "val2",
		})
		infos := GetInfo(rootErr)

		require.Len(t, infos, 1)
		assert.Equal(t, Fields{
			"key1": "val1",
			"key2": "val2",
		}, infos[0].Fields)
	})
}

func Test_W(t *testing.T) {
	t.Run("W fills message based on the function it was used in", func(t *testing.T) {
		err := dummyFunc()
		assert.Equal(t, "dummyFunc: rootErr", err.Error())
	})
}

// NCError.Unwrap supports Unwrap interface
func Test_NCError_Unwrap(t *testing.T) {
	wrappedErr := Wrap(rootErr, "message", nil)

	assert.ErrorIs(t, wrappedErr, rootErr)
}
