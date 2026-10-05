package app

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
)

// UpgradeName is the name of the software upgrade handled by this binary.
// It must match the plan name of the governance software-upgrade proposal.
const UpgradeName = "v0.1.3"

// setupUpgradeHandlers registers the handlers of the software upgrades supported by this binary.
func (app *App) setupUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// this upgrade carries no state changes: it only runs the pending module migrations
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)
}
