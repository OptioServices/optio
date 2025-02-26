package keeper

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"cosmossdk.io/store/prefix"
	"github.com/OptioServices/optio/x/distro/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DailyDistributionTotalAll(ctx context.Context, req *types.QueryAllDailyDistributionTotalRequest) (*types.QueryAllDailyDistributionTotalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	dailyDistributionTotals := make(map[string]uint64)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	dailyDistributionTotalStore := prefix.NewStore(store, types.KeyPrefix(types.DailyDistributionTotalKeyPrefix))

	pageRes, err := query.Paginate(dailyDistributionTotalStore, req.Pagination, func(key []byte, value []byte) error {
		val, err := strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return nil
		}

		dailyDistributionTotals[string(key)] = val
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	params := k.GetParams(ctx)

	totals := make([]types.DailyDistributionTotalEntry, 0, len(dailyDistributionTotals))
	for date, total := range dailyDistributionTotals {
		totals = append(totals, types.DailyDistributionTotalEntry{
			Date:   date,
			Amount: fmt.Sprintf("%d%s", total, params.Denom),
		})
	}

	sort.Slice(totals, func(i, j int) bool {
		return totals[i].Date > totals[j].Date
	})

	return &types.QueryAllDailyDistributionTotalResponse{DailyDistributionTotals: totals, Pagination: pageRes}, nil
}

func (k Keeper) DailyDistributionTotal(ctx context.Context, req *types.QueryGetDailyDistributionTotalRequest) (*types.QueryGetDailyDistributionTotalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetDailyDistributionTotal(
		ctx,
		req.Date,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDailyDistributionTotalResponse{Total: val}, nil
}
