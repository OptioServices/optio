package v2

import (
	store "cosmossdk.io/store/types"
	"github.com/OptioServices/optio/app/upgrades"
	optiotypes "github.com/OptioServices/optio/x/distribute/types"
)

const UpgradeName = "v2"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{optiotypes.ModuleName},
	},
}
