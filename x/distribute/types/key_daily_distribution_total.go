package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// DailyDistributionTotalKeyPrefix is the prefix to retrieve all DailyDistributionTotal
	DailyDistributionTotalKeyPrefix = "distributed/"
)

// DailyDistributionTotalKey returns the store key to retrieve a DailyDistributionTotal from the index fields
func DailyDistributionTotalKey(
	date string,
) []byte {
	var key []byte

	dateBytes := []byte(date)
	key = append(key, dateBytes...)

	return key
}
