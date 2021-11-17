// this package contains related to stack trace, thus are susceptible to changes in source coude.
// there are no anonymous functions here since they screw up with stack trace
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
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo_NewErrorReturnsFuncName:15",
		infos[0].FuncName,
	)
}

func Test_GetInfo_WrappedErrorReturnsFuncName(t *testing.T) {
	rootErr := errors.Wrap(fmt.Errorf("rootErr"), "wrappedErr", nil)
	infos := errors.GetInfo(rootErr)

	require.Len(t, infos, 2)
	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo_WrappedErrorReturnsFuncName:26",
		infos[0].FuncName,
	)
}
