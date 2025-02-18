package keeper

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/OptioServices/optio/x/distribute/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// SetDailyDistributionTotal set a specific dailyDistributionTotal in the store from its index
func (k Keeper) SetDailyDistributionTotal(ctx context.Context, dailyDistributionTotal uint64, date string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyDistributionTotalKeyPrefix))
	store.Set(types.DailyDistributionTotalKey(
		date,
	), []byte(strconv.FormatUint(dailyDistributionTotal, 10)))
}

// GetDailyDistributionTotal returns a dailyDistributionTotal from its index
func (k Keeper) GetDailyDistributionTotal(
	ctx context.Context,
	date string,

) (val uint64, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyDistributionTotalKeyPrefix))

	b := store.Get(types.DailyDistributionTotalKey(
		date,
	))
	if b == nil {
		return val, false
	}

	val, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return val, false
	}
	return val, true
}

// RemoveDailyDistributionTotal removes a dailyDistributionTotal from the store
func (k Keeper) RemoveDailyDistributionTotal(
	ctx context.Context,
	date string,

) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyDistributionTotalKeyPrefix))
	store.Delete(types.DailyDistributionTotalKey(
		date,
	))
}

// GetAllDailyDistributionTotal returns all dailyDistributionTotal
func (k Keeper) GetAllDailyDistributionTotal(ctx context.Context) map[string]uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.DailyDistributionTotalKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	totalMap := make(map[string]uint64)

	for ; iterator.Valid(); iterator.Next() {
		value := iterator.Value()
		if len(value) != 8 {
			panic(fmt.Errorf("invalid value length for uint64: %d", len(value)))
		}
		val := binary.BigEndian.Uint64(value)
		totalMap[string(iterator.Key())] = val
	}

	return totalMap
}

// UpdateDailyDistributionTotal updates the total for a specific date
func (k Keeper) UpdateDailyDistributionTotal(ctx context.Context, date string, amount uint64) error {
	dailyTotal, found := k.GetDailyDistributionTotal(ctx, date)
	if found {
		k.SetDailyDistributionTotal(ctx, dailyTotal+amount, date)
	} else {
		k.SetDailyDistributionTotal(ctx, amount, date)
	}

	return nil
}
