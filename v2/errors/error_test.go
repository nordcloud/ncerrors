package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyError struct {
	Msg string
}

func (e dummyError) Error() string {
	return e.Msg
}

func dummyFunc() error {
	return W(dummyError{Msg: "rootError"})
}

func Test_Wrap(t *testing.T) {
	t.Run("Wrapped error returns concatenation of messages", func(t *testing.T) {
		rootErr := dummyError{"rootErr"}
		wrappedErr := Wrap(rootErr, "message", nil)

		assert.Equal(t, "message: rootErr", wrappedErr.Error())
	})

	t.Run("Wrapped error is unwrappable", func(t *testing.T) {
		rootErr := dummyError{"rootErr"}
		wrappedErr := Wrap(rootErr, "message", nil)

		assert.Equal(t, rootErr, Unwrap(wrappedErr))
	})
}

func Test_GetInfo(t *testing.T) {
	t.Run("Sentinel errors returns only message", func(t *testing.T) {
		rootErr := dummyError{"rootError"}
		infos := GetInfo(rootErr)

		require.Len(t, infos, 1)
		assert.Equal(t, "rootError", infos[0].Message)
		assert.Nil(t, infos[0].Fields)
		assert.Nil(t, infos[0].StackTrace)
	})

	t.Run("NCError errors returns message", func(t *testing.T) {
		rootErr := New("rootError", nil)
		infos := GetInfo(rootErr)

		require.Len(t, infos, 1)
		assert.Equal(t, "rootError", infos[0].Message)
	})

	t.Run("NCError errors returns fields", func(t *testing.T) {
		rootErr := New("rootError", Fields{
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
		assert.Equal(t, "dummyFunc: rootError", err.Error())
	})
}

func Test_Is(t *testing.T) {
	rootErr := dummyError{"rootErr"}
	wrappedErr := Wrap(rootErr, "message", nil)

	assert.ErrorIs(t, wrappedErr, rootErr)
}

func Test_As(t *testing.T) {
	rootErr := dummyError{"rootErr"}
	wrappedErr := Wrap(rootErr, "message", nil)

	var dummyErr dummyError

	assert.ErrorAs(t, wrappedErr, &dummyErr)
}
