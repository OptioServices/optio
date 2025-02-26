package keeper

import (
	"github.com/OptioServices/optio/x/distro/types"
)

var _ types.QueryServer = Keeper{}
