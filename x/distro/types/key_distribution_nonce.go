package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// DistributionNonceKeyPrefix is the prefix to retrieve all DistributionNonce
	DistributionNonceKeyPrefix = "nonce/"
)

// DistributionNonceKey returns the store key to retrieve a DistributionNonce from the index fields
func DistributionNonceKey(
	nonce string,
) []byte {
	var key []byte

	nonceBytes := []byte(nonce)
	key = append(key, nonceBytes...)

	return key
}
