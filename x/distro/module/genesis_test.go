package optio_test

import (
	"testing"

	keepertest "github.com/OptioServices/optio/testutil/keeper"
	"github.com/OptioServices/optio/testutil/nullify"
	distro "github.com/OptioServices/optio/x/distro/module"
	"github.com/OptioServices/optio/x/distro/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.DistributeKeeper(t)
	distro.InitGenesis(ctx, k, genesisState)
	got := distro.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
