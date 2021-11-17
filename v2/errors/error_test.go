package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var rootSentinelErr = fmt.Errorf("rootSentinelErr")

func dummyFunc() error {
	return W(rootSentinelErr)
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

	t.Run("Wrapped NCError", func(t *testing.T) {
	})
}

func Test_W(t *testing.T) {
	t.Run("W fills message based on the function it was used in", func(t *testing.T) {
		err := dummyFunc()
		assert.Equal(t, "dummyFunc: rootSentinelErr", err.Error())
	})
}
