// this package contains related to stack trace, thus are susceptible to changes in source coude.
// Files in this package should be kept relatively small in order to avoid breaking changes in multiple tests with a small change.
// there are no anonymous functions here since they mess with stack trace
package stacktest

import (
	"testing"

	"github.com/nordcloud/ncerrors/v2/errors"
	"github.com/stretchr/testify/assert"
)

func Test_GetInfo_NewErrorReturnsStackTrace(t *testing.T) {
	rootErr := errors.New("rootError", nil)
	infos := errors.GetInfo(rootErr)

	assert.GreaterOrEqual(t, 1, len(infos))
	assert.Equal(t,
		"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo_NewErrorReturnsStackTrace:14",
		infos[0].StackTrace[0],
	)
}
