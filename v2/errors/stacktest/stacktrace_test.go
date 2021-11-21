// this package contains related to stack trace, thus are susceptible to changes in source coude.
// Files in this package should be kept relatively small in order to avoid breaking changes in multiple tests with a small change.
// there are no anonymous functions here since they mess with stack trace
package stacktest

import (
	"fmt"
	"testing"

	"github.com/nordcloud/ncerrors/v2/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newNCErr() error {
	return errors.New("rootNCError", nil)
}

func Test_GetInfo_NewErrorReturnsStackTrace(t *testing.T) {
	infos := errors.GetInfo(newNCErr())

	require.GreaterOrEqual(t, len(infos), 1)
	require.GreaterOrEqual(t, len(infos[0].StackTrace), 2)
	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.newNCErr:16",
		infos[0].StackTrace[0],
	)

	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo_NewErrorReturnsStackTrace:20",
		infos[0].StackTrace[1],
	)
}

func wrapSentinelErr() error {
	return errors.Wrap(fmt.Errorf("rootSentinelErr"), "", nil)
}

func Test_GetInfo_WrapReturnsStackTrace(t *testing.T) {
	infos := errors.GetInfo(wrapSentinelErr())

	require.GreaterOrEqual(t, len(infos), 2)
	require.GreaterOrEqual(t, len(infos[0].StackTrace), 2)
	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.wrapSentinelErr:36",
		infos[0].StackTrace[0],
	)

	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo_WrapReturnsStackTrace:40",
		infos[0].StackTrace[1],
	)
}

func wSentinelErr() error {
	return errors.W(fmt.Errorf("rootSentinelErr"))
}

func Test_GetInfo_WReturnsStackTrace(t *testing.T) {
	infos := errors.GetInfo(wSentinelErr())

	require.GreaterOrEqual(t, len(infos), 1)
	require.GreaterOrEqual(t, len(infos[0].StackTrace), 2)
	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.wSentinelErr:56",
		infos[0].StackTrace[0],
	)

	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo_WReturnsStackTrace:60",
		infos[0].StackTrace[1],
	)
}
