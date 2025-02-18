package keeper

import (
	"context"
	"encoding/binary"
	"fmt"

	"cosmossdk.io/store/prefix"
	"github.com/OptioServices/optio/x/distribute/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DailyDistributionTotalAll(ctx context.Context, req *types.QueryAllDailyDistributionTotalRequest) (*types.QueryAllDailyDistributionTotalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var dailyDistributionTotals map[string]uint64

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	dailyDistributionTotalStore := prefix.NewStore(store, types.KeyPrefix(types.DailyDistributionTotalKeyPrefix))

	pageRes, err := query.Paginate(dailyDistributionTotalStore, req.Pagination, func(key []byte, value []byte) error {
		if len(value) != 8 {
			return fmt.Errorf("invalid value length for uint64: %d", len(value))
		}
		val := binary.BigEndian.Uint64(value)

		dailyDistributionTotals[string(key)] = val
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDailyDistributionTotalResponse{DailyDistributionTotals: dailyDistributionTotals, Pagination: pageRes}, nil
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
