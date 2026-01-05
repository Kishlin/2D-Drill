package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type UpgradeSystem struct {
	engineShop   *entities.UpgradeShop
	hullShop     *entities.UpgradeShop
	fuelTankShop *entities.UpgradeShop
}

func NewUpgradeSystem(
	engineShop, hullShop, fuelTankShop *entities.UpgradeShop,
) *UpgradeSystem {
	return &UpgradeSystem{
		engineShop:   engineShop,
		hullShop:     hullShop,
		fuelTankShop: fuelTankShop,
	}
}

func (us *UpgradeSystem) ProcessUpgrade(
	player *entities.Player,
	inputState input.InputState,
) {
	if !inputState.Sell {
		return
	}

	// Check each shop and attempt upgrade (only one can be in range at a time)
	if us.engineShop.IsPlayerInRange(player) {
		us.tryUpgradeEngine(player)
		return
	}
	if us.hullShop.IsPlayerInRange(player) {
		us.tryUpgradeHull(player)
		return
	}
	if us.fuelTankShop.IsPlayerInRange(player) {
		us.tryUpgradeFuelTank(player)
		return
	}
}

func (us *UpgradeSystem) tryUpgradeEngine(player *entities.Player) {
	cost := entities.GetEngineNextCost(player.Upgrades.Engine)
	if cost == 0 {
		return // Already at max level
	}

	if player.Money < cost {
		return // Cannot afford
	}

	player.Money -= cost
	player.Upgrades.Engine++
}

func (us *UpgradeSystem) tryUpgradeHull(player *entities.Player) {
	cost := entities.GetHullNextCost(player.Upgrades.Hull)
	if cost == 0 {
		return // Already at max level
	}

	if player.Money < cost {
		return // Cannot afford
	}

	player.Money -= cost
	player.Upgrades.Hull++
}

func (us *UpgradeSystem) tryUpgradeFuelTank(player *entities.Player) {
	cost := entities.GetFuelTankNextCost(player.Upgrades.FuelTank)
	if cost == 0 {
		return // Already at max level
	}

	if player.Money < cost {
		return // Cannot afford
	}

	player.Money -= cost
	player.Upgrades.FuelTank++
}

func (us *UpgradeSystem) GetEngineShop() *entities.UpgradeShop {
	return us.engineShop
}

func (us *UpgradeSystem) GetHullShop() *entities.UpgradeShop {
	return us.hullShop
}

func (us *UpgradeSystem) GetFuelTankShop() *entities.UpgradeShop {
	return us.fuelTankShop
}
