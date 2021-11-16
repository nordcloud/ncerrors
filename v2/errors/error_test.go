package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
