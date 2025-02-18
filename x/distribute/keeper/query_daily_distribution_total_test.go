package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/OptioServices/optio/testutil/keeper"
	"github.com/OptioServices/optio/testutil/nullify"
	"github.com/OptioServices/optio/x/distribute/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestDailyDistributionTotalQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.DistributeKeeper(t)
	msgs := createNDailyDistributionTotal(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetDailyDistributionTotalRequest
		response *types.QueryGetDailyDistributionTotalResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetDailyDistributionTotalRequest{
				Date: "2025-03-01",
			},
			response: &types.QueryGetDailyDistributionTotalResponse{Total: msgs["2025-03-01"]},
		},
		{
			desc: "Second",
			request: &types.QueryGetDailyDistributionTotalRequest{
				Date: "2025-03-02",
			},
			response: &types.QueryGetDailyDistributionTotalResponse{Total: msgs["2025-03-02"]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDailyDistributionTotalRequest{
				Date: "2025-03-03",
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.DailyDistributionTotal(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestDailyDistributionTotalQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.DistributeKeeper(t)
	msgs := createNDailyDistributionTotal(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllDailyDistributionTotalRequest {
		return &types.QueryAllDailyDistributionTotalRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.DailyDistributionTotalAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DailyDistributionTotals), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.DailyDistributionTotals),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.DailyDistributionTotalAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DailyDistributionTotals), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.DailyDistributionTotals),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.DailyDistributionTotalAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.DailyDistributionTotals),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.DailyDistributionTotalAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
