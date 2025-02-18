package optio

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/OptioServices/optio/x/distribute/keeper"
	"github.com/OptioServices/optio/x/distribute/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the dailyDistributionTotal
	for date, elem := range genState.DailyDistributionTotals {
		k.SetDailyDistributionTotal(ctx, elem, date)
	}
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
