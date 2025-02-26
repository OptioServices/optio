package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/OptioServices/optio/x/distro/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetDailyDistributionTotal set a specific dailyDistributionTotal in the store from its index
func (k Keeper) SetDailyDistributionTotal(ctx context.Context, date string, total uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.DailyDistributionTotalKeyPrefix))
	store.Set([]byte(date), []byte(strconv.FormatUint(total, 10)))
}

// GetDailyDistributionTotal returns a dailyDistributionTotal from its index
func (k Keeper) GetDailyDistributionTotal(ctx context.Context, date string) (uint64, bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.DailyDistributionTotalKeyPrefix))
	b := store.Get([]byte(date))
	if b == nil {
		return 0, false
	}
	val, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0, false
	}
	return val, true
}

func (k Keeper) GetDailyDistributionTotals(ctx context.Context) map[string]uint64 {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.DailyDistributionTotalKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	totals := make(map[string]uint64)
	for ; iterator.Valid(); iterator.Next() {
		date := string(iterator.Key())
		if val, err := strconv.ParseUint(string(iterator.Value()), 10, 64); err == nil {
			totals[date] = val
		}
	}
	return totals
}
