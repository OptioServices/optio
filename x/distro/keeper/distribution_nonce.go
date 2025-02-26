package keeper

import (
	"context"
	"strconv"

	"cosmossdk.io/store/prefix"
	"github.com/OptioServices/optio/x/distro/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// HasNonceBeenUsed checks if a nonce has been used before.
func (k Keeper) HasNonceBeenUsed(ctx context.Context, nonce uint64) bool {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.DistributionNonceKeyPrefix))
	return store.Has([]byte(strconv.FormatUint(nonce, 10)))
}

// SetNonceUsed marks a nonce as used.
func (k Keeper) SetNonceUsed(ctx context.Context, nonce uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.DistributionNonceKeyPrefix))
	store.Set([]byte(strconv.FormatUint(nonce, 10)), []byte{1}) // Value is irrelevant, just presence matters
}
