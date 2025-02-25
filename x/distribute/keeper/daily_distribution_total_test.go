package keeper_test

import (
	"context"
	"strconv"
	"testing"

	keepertest "github.com/OptioServices/optio/testutil/keeper"
	"github.com/OptioServices/optio/testutil/nullify"
	"github.com/OptioServices/optio/x/distribute/keeper"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNDailyDistributionTotal(keeper keeper.Keeper, ctx context.Context, n int) map[string]uint64 {
	items := make(map[string]uint64)
	for i := 0; i < n; i++ {
		items[strconv.Itoa(i)] = uint64(i)

		keeper.SetDailyDistributionTotal(ctx, strconv.Itoa(i), uint64(i))
	}
	return items
}

func TestDailyDistributionTotalGet(t *testing.T) {
	keeper, ctx := keepertest.DistributeKeeper(t)
	items := createNDailyDistributionTotal(keeper, ctx, 10)
	for i, item := range items {
		rst, found := keeper.GetDailyDistributionTotal(ctx,
			i,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func TestDailyDistributionTotalGetAll(t *testing.T) {
	keeper, ctx := keepertest.DistributeKeeper(t)
	items := createNDailyDistributionTotal(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetDailyDistributionTotals(ctx)),
	)
}
