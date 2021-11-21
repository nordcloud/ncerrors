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

func Test_GetInfo_NewErrorReturnsFuncName(t *testing.T) {
	rootErr := errors.New("rootError", nil)
	infos := errors.GetInfo(rootErr)

	require.Len(t, infos, 1)
	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo_NewErrorReturnsFuncName:16",
		infos[0].FuncName,
	)
}

func Test_GetInfo_WrappedErrorReturnsFuncName(t *testing.T) {
	rootErr := errors.Wrap(fmt.Errorf("rootErr"), "wrappedErr", nil)
	infos := errors.GetInfo(rootErr)

	require.Len(t, infos, 2)
	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo_WrappedErrorReturnsFuncName:27",
		infos[0].FuncName,
	)
}

func Test_GetInfo_WErrorReturnsFuncName(t *testing.T) {
	rootErr := errors.W(fmt.Errorf("rootErr"))
	infos := errors.GetInfo(rootErr)

	require.Len(t, infos, 2)
	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo_WErrorReturnsFuncName:38",
		infos[0].FuncName,
	)
}
