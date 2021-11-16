// this package contains related to stack trace, thus are susceptible to changes in source coude.
package stacktest

import (
	"testing"

	"github.com/nordcloud/ncerrors/v2/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetInfo(t *testing.T) {
	t.Run("NCError errors returns stack trace", func(t *testing.T) {
		rootErr := errors.New("rootError", nil)
		infos := errors.GetInfo(rootErr)

		require.Len(t, infos, 1)
		assert.Equal(t,
			"github.com/nordcloud/ncerrors/v2/errors/stacktest.Test_GetInfo.func1:14",
			infos[0].FuncName,
		)
	})
}
