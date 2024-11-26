package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/OptioServices/optio/testutil/keeper"
	"github.com/OptioServices/optio/x/optio/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.OptioKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
