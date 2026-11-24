package app

import (
	"context"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpgradeName is the name of the software upgrade handled by this binary.
// It must match the plan name of the governance software-upgrade proposal.
//
// v0.3.0 moves the chain from Cosmos SDK v0.53 / CometBFT v0.38 to
// Cosmos SDK v0.55 / CometBFT v0.40 in a single coordinated upgrade.
const UpgradeName = "v0.3.0"

// storeUpgrades lists the store changes applied at the upgrade height.
//
// The deleted stores belong to modules this binary no longer mounts:
//   - params: x/params was removed from the SDK in v0.55; every module manages
//     its own parameters since v0.47 and the store only held the legacy copy.
//   - group: x/group moved to Cosmos Enterprise (commercial license).
//   - circuit, nft: moved to the unmaintained contrib tree of the SDK.
//
// All four stores are empty on worrell-testnet-2 (checked on 2026-11-24).
var storeUpgrades = storetypes.StoreUpgrades{
	Added:   []string{},
	Deleted: []string{"params", "group", "circuit", "nft"},
}

// setupUpgradeHandlers registers the handlers of the software upgrades supported by this binary
// and, when the node is restarting at the upgrade height, the store loader that applies the
// store changes before the module migrations run.
func (app *App) setupUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// the module migrations carry the state changes of this upgrade
			// (auth, staking with the key_rotation_fee param, feegrant, mint)
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// configure the store loader that checks if version == upgradeHeight and applies the store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
