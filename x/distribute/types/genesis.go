package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		DailyDistributionTotals: map[string]uint64{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in dailyDistributionTotal
	dailyDistributionTotalIndexMap := make(map[string]struct{})

	for date := range gs.DailyDistributionTotals {
		index := string(DailyDistributionTotalKey(date))
		if _, ok := dailyDistributionTotalIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for dailyDistributionTotal")
		}
		dailyDistributionTotalIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
