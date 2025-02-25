package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/OptioServices/optio/x/distribute/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// HasNonceBeenUsed checks if a nonce has been used before.
func (k Keeper) HasNonceBeenUsed(ctx context.Context, nonce string) bool {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.DistributionNonceKeyPrefix))
	return store.Has([]byte(nonce))
}

// SetNonceUsed marks a nonce as used.
func (k Keeper) SetNonceUsed(ctx context.Context, nonce string) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.DistributionNonceKeyPrefix))
	store.Set([]byte(nonce), []byte{1}) // Value is irrelevant, just presence matters
}
