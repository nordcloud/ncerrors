package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	rootSentinelErr = fmt.Errorf("rootSentinelErr")
	rootNCErr       = New("rootNCErr", nil)
)

func wSentinelErr() error {
	return W(rootSentinelErr)
}

func w2SentinelErr() error {
	return W(wSentinelErr())
}

func wNCErr() error {
	return W(rootNCErr)
}

func w2NCErr() error {
	return W(wNCErr())
}

// A custom error that retains all properties of NCError, like stacktrace, fields, func name etc.
type customErr struct {
	NCError
}

func (e customErr) Custom() bool {
	return true
}

func Test_New(t *testing.T) {
	t.Run("New NCError returns error message", func(t *testing.T) {
		rootNCErr := New("rootNCErr", nil)

		assert.Equal(t, "rootNCErr", rootNCErr.Error())
	})
}

func Test_Wrap(t *testing.T) {
	t.Run("Wrap2 -> Wrap1 -> Root NCError returns concatenation of messages", func(t *testing.T) {
		rootNCErr := New("rootNCErr", nil)
		wrap1Err := Wrap(rootNCErr, "wrap1", nil)
		wrap2Err := Wrap(wrap1Err, "wrap2", nil)

		assert.Equal(t, "wrap2: wrap1: rootNCErr", wrap2Err.Error())
	})

	t.Run("Wrap2 -> Wrap1 -> Root sentinel returns concatenation of messages", func(t *testing.T) {
		wrap1Err := Wrap(rootSentinelErr, "wrap1", nil)
		wrap2Err := Wrap(wrap1Err, "wrap2", nil)

		assert.Equal(t, "wrap2: wrap1: rootSentinelErr", wrap2Err.Error())
	})

	t.Run("Wrapped error is unwrappable", func(t *testing.T) {
		wrappedErr := Wrap(rootSentinelErr, "message", nil)

		assert.Equal(t, rootSentinelErr, Unwrap(wrappedErr))
	})
}

func Test_GetInfo(t *testing.T) {
	t.Run("Wrap2 -> Wrap1 -> Root NCError", func(t *testing.T) {
		t.Run("returns message", func(t *testing.T) {
			rootNCErr := New("rootNCErr", nil)
			wrap1Err := Wrap(rootNCErr, "wrap1", nil)
			wrap2Err := Wrap(wrap1Err, "wrap2", nil)

			infos := GetInfo(wrap2Err)

			require.Len(t, infos, 3)
			assert.Equal(t, "wrap2", infos[0].Message)
			assert.Equal(t, "wrap1", infos[1].Message)
			assert.Equal(t, "rootNCErr", infos[2].Message)
		})

		t.Run("returns fields", func(t *testing.T) {
			rootNCErr := New(
				"rootNCErr",
				Fields{
					"key1": "val1",
				},
			)
			wrap1Err := Wrap(rootNCErr, "wrap1",
				Fields{
					"key2": "val2",
				},
			)
			wrap2Err := Wrap(wrap1Err, "wrap2",
				Fields{
					"key3": "val3",
				},
			)

			infos := GetInfo(wrap2Err)

			require.Len(t, infos, 3)
			assert.Equal(t,
				Fields{
					"key3": "val3",
				},
				infos[0].Fields,
			)
			assert.Equal(t,
				Fields{
					"key2": "val2",
				},
				infos[1].Fields,
			)
			assert.Equal(t,
				Fields{
					"key1": "val1",
				},
				infos[2].Fields,
			)
		})

		t.Run("returns non empty stack trace for root error", func(t *testing.T) {
			rootNCErr := New("rootNCErr", nil)
			wrap1Err := Wrap(rootNCErr, "wrap1", nil)
			wrap2Err := Wrap(wrap1Err, "wrap2", nil)

			infos := GetInfo(wrap2Err)

			require.Len(t, infos, 3)
			assert.Empty(t, infos[0].StackTrace)
			assert.Empty(t, infos[1].StackTrace)
			assert.NotEmpty(t, infos[2].StackTrace)
		})

		t.Run("returns non empty func name for all errors", func(t *testing.T) {
			rootNCErr := New("rootNCErr", nil)
			wrap1Err := Wrap(rootNCErr, "wrap1", nil)
			wrap2Err := Wrap(wrap1Err, "wrap2", nil)

			infos := GetInfo(wrap2Err)

			require.Len(t, infos, 3)
			assert.NotEmpty(t, infos[0].FuncName)
			assert.NotEmpty(t, infos[1].FuncName)
			assert.NotEmpty(t, infos[2].FuncName)
		})
	})

	t.Run("Wrap2 -> Wrap1 -> Root sentinel error", func(t *testing.T) {
		t.Run("returns message", func(t *testing.T) {
			wrap1Err := Wrap(rootSentinelErr, "wrap1", nil)
			wrap2Err := Wrap(wrap1Err, "wrap2", nil)

			infos := GetInfo(wrap2Err)

			require.Len(t, infos, 3)
			assert.Equal(t, "wrap2", infos[0].Message)
			assert.Equal(t, "wrap1", infos[1].Message)
			assert.Equal(t, "rootSentinelErr", infos[2].Message)
		})

		t.Run("returns fields for NC errors", func(t *testing.T) {
			wrap1Err := Wrap(rootSentinelErr, "wrap1",
				Fields{
					"key2": "val2",
				},
			)
			wrap2Err := Wrap(wrap1Err, "wrap2",
				Fields{
					"key3": "val3",
				},
			)

			infos := GetInfo(wrap2Err)

			require.Len(t, infos, 3)
			assert.Equal(t,
				Fields{
					"key3": "val3",
				},
				infos[0].Fields,
			)
			assert.Equal(t,
				Fields{
					"key2": "val2",
				},
				infos[1].Fields,
			)
			assert.Empty(t, infos[2].Fields)
		})

		t.Run("returns non empty stack trace for first NC error", func(t *testing.T) {
			wrap1Err := Wrap(rootSentinelErr, "wrap1", nil)
			wrap2Err := Wrap(wrap1Err, "wrap2", nil)

			infos := GetInfo(wrap2Err)

			require.Len(t, infos, 3)
			assert.Empty(t, infos[0].StackTrace)
			assert.NotEmpty(t, infos[1].StackTrace)
			assert.Empty(t, infos[2].StackTrace)
		})

		t.Run("returns non empty func name for all NC errors", func(t *testing.T) {
			wrap1Err := Wrap(rootSentinelErr, "wrap1", nil)
			wrap2Err := Wrap(wrap1Err, "wrap2", nil)

			infos := GetInfo(wrap2Err)

			require.Len(t, infos, 3)
			assert.NotEmpty(t, infos[0].FuncName)
			assert.NotEmpty(t, infos[1].FuncName)
			assert.Empty(t, infos[2].FuncName)
		})
	})

	t.Run("W2 -> W1 -> Root sentinel error", func(t *testing.T) {
		t.Run("returns message", func(t *testing.T) {
			w2Err := w2SentinelErr()

			infos := GetInfo(w2Err)

			require.Len(t, infos, 3)
			assert.Equal(t, "w2SentinelErr", infos[0].Message)
			assert.Equal(t, "wSentinelErr", infos[1].Message)
			assert.Equal(t, "rootSentinelErr", infos[2].Message)
		})

		t.Run("returns non empty stack trace for first NC error", func(t *testing.T) {
			w2Err := w2SentinelErr()

			infos := GetInfo(w2Err)

			require.Len(t, infos, 3)
			assert.Empty(t, infos[0].StackTrace)
			assert.NotEmpty(t, infos[1].StackTrace)
			assert.Empty(t, infos[2].StackTrace)
		})

		t.Run("returns non empty func name for all NC errors", func(t *testing.T) {
			w2Err := w2SentinelErr()

			infos := GetInfo(w2Err)

			require.Len(t, infos, 3)
			assert.NotEmpty(t, infos[0].FuncName)
			assert.NotEmpty(t, infos[1].FuncName)
			assert.Empty(t, infos[2].FuncName)
		})
	})

	t.Run("Wrap -> NewWithErr -> Root NCError", func(t *testing.T) {
		t.Run("returns message", func(t *testing.T) {
			rootNCErr := New("rootNCErr", nil)
			newWithErr := customErr{NCError: NewWithErr(rootNCErr, "newWithErr", nil)}
			wrapErr := Wrap(newWithErr, "wrap", nil)

			infos := GetInfo(wrapErr)

			require.Len(t, infos, 3)
			assert.Equal(t, "wrap", infos[0].Message)
			assert.Equal(t, "newWithErr", infos[1].Message)
			assert.Equal(t, "rootNCErr", infos[2].Message)
		})

		t.Run("returns fields", func(t *testing.T) {
			rootNCErr := New(
				"rootNCErr",
				Fields{
					"key1": "val1",
				},
			)
			newWithErr := customErr{NCError: NewWithErr(rootNCErr, "newWithErr",
				Fields{
					"key2": "val2",
				},
			)}
			wrapErr := Wrap(newWithErr, "wrap",
				Fields{
					"key3": "val3",
				},
			)

			infos := GetInfo(wrapErr)

			require.Len(t, infos, 3)
			assert.Equal(t,
				Fields{
					"key3": "val3",
				},
				infos[0].Fields,
			)
			assert.Equal(t,
				Fields{
					"key2": "val2",
				},
				infos[1].Fields,
			)
			assert.Equal(t,
				Fields{
					"key1": "val1",
				},
				infos[2].Fields,
			)
		})

		t.Run("returns non empty stack trace for root error", func(t *testing.T) {
			rootNCErr := New("rootNCErr", nil)
			newWithErr := customErr{NCError: NewWithErr(rootNCErr, "newWithErr", nil)}
			wrapErr := Wrap(newWithErr, "wrap", nil)

			infos := GetInfo(wrapErr)

			require.Len(t, infos, 3)
			assert.Empty(t, infos[0].StackTrace)
			assert.Empty(t, infos[1].StackTrace)
			assert.NotEmpty(t, infos[2].StackTrace)
		})

		t.Run("returns non empty func name for all errors", func(t *testing.T) {
			rootNCErr := New("rootNCErr", nil)
			newWithErr := customErr{NCError: NewWithErr(rootNCErr, "newWithErr", nil)}
			wrapErr := Wrap(newWithErr, "wrap", nil)

			infos := GetInfo(wrapErr)

			require.Len(t, infos, 3)
			assert.NotEmpty(t, infos[0].FuncName)
			assert.NotEmpty(t, infos[1].FuncName)
			assert.NotEmpty(t, infos[2].FuncName)
		})
	})

	t.Run("Wrap -> NewWithErr -> Root sentinel error", func(t *testing.T) {
		t.Run("returns message", func(t *testing.T) {
			newWithErr := customErr{NCError: NewWithErr(rootSentinelErr, "newWithErr", nil)}
			wrapErr := Wrap(newWithErr, "wrap", nil)

			infos := GetInfo(wrapErr)

			require.Len(t, infos, 3)
			assert.Equal(t, "wrap", infos[0].Message)
			assert.Equal(t, "newWithErr", infos[1].Message)
			assert.Equal(t, "rootSentinelErr", infos[2].Message)
		})

		t.Run("returns fields", func(t *testing.T) {
			newWithErr := customErr{NCError: NewWithErr(rootSentinelErr, "newWithErr",
				Fields{
					"key2": "val2",
				},
			)}
			wrapErr := Wrap(newWithErr, "wrap",
				Fields{
					"key3": "val3",
				},
			)

			infos := GetInfo(wrapErr)

			require.Len(t, infos, 3)
			assert.Equal(t,
				Fields{
					"key3": "val3",
				},
				infos[0].Fields,
			)
			assert.Equal(t,
				Fields{
					"key2": "val2",
				},
				infos[1].Fields,
			)
			assert.Nil(t, infos[2].Fields)
		})

		t.Run("returns non empty stack trace for root error", func(t *testing.T) {
			newWithErr := customErr{NCError: NewWithErr(rootSentinelErr, "newWithErr", nil)}
			wrapErr := Wrap(newWithErr, "wrap", nil)

			infos := GetInfo(wrapErr)

			require.Len(t, infos, 3)
			assert.Empty(t, infos[0].StackTrace)
			assert.NotEmpty(t, infos[1].StackTrace)
			assert.Empty(t, infos[2].StackTrace)
		})

		t.Run("returns non empty func name for all nc errors", func(t *testing.T) {
			newWithErr := customErr{NCError: NewWithErr(rootSentinelErr, "newWithErr", nil)}
			wrapErr := Wrap(newWithErr, "wrap", nil)

			infos := GetInfo(wrapErr)

			require.Len(t, infos, 3)
			assert.NotEmpty(t, infos[0].FuncName)
			assert.NotEmpty(t, infos[1].FuncName)
			assert.Empty(t, infos[2].FuncName)
		})
	})
}

func Test_W(t *testing.T) {
	t.Run("W2 -> W1 -> Root sentinel error returns concatenation of messages", func(t *testing.T) {
		assert.Equal(t, "w2SentinelErr: wSentinelErr: rootSentinelErr", w2SentinelErr().Error())
	})

	t.Run("W2 -> W1 -> Root NCError returns concatenation of messages", func(t *testing.T) {
		assert.Equal(t, "w2NCErr: wNCErr: rootNCErr", w2NCErr().Error())
	})

	t.Run("Wrapped error is unwrappable", func(t *testing.T) {
		wrappedErr := W(rootSentinelErr)

		assert.Equal(t, rootSentinelErr, Unwrap(wrappedErr))
	})
}
