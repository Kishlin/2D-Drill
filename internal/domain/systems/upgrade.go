package systems

import (
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/input"
)

type UpgradeSystem struct {
	engineShop     *entities.EngineUpgradeShop
	hullShop       *entities.HullUpgradeShop
	fuelTankShop   *entities.FuelTankUpgradeShop
	cargoHoldShop  *entities.CargoHoldUpgradeShop
	heatShieldShop *entities.HeatShieldUpgradeShop
}

func NewUpgradeSystem(
	engineShop *entities.EngineUpgradeShop,
	hullShop *entities.HullUpgradeShop,
	fuelTankShop *entities.FuelTankUpgradeShop,
	cargoHoldShop *entities.CargoHoldUpgradeShop,
	heatShieldShop *entities.HeatShieldUpgradeShop,
) *UpgradeSystem {
	return &UpgradeSystem{
		engineShop:     engineShop,
		hullShop:       hullShop,
		fuelTankShop:   fuelTankShop,
		cargoHoldShop:  cargoHoldShop,
		heatShieldShop: heatShieldShop,
	}
}

func (us *UpgradeSystem) ProcessUpgrade(
	player *entities.Player,
	inputState input.InputState,
) {
	if !inputState.Sell {
		return
	}

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
	if us.cargoHoldShop.IsPlayerInRange(player) {
		us.tryUpgradeCargoHold(player)
		return
	}
	if us.heatShieldShop.IsPlayerInRange(player) {
		us.tryUpgradeHeatShield(player)
		return
	}
}

func (us *UpgradeSystem) tryUpgradeEngine(player *entities.Player) {
	entry := us.engineShop.GetNextEngine(player.Engine.Tier())
	if entry == nil {
		return // Max level reached
	}

	if !player.CanAfford(entry.Price) {
		return
	}

	player.BuyEngine(entry.Engine, entry.Price)
}

func (us *UpgradeSystem) tryUpgradeHull(player *entities.Player) {
	entry := us.hullShop.GetNextHull(player.Hull.Tier())
	if entry == nil {
		return // Max level reached
	}

	if !player.CanAfford(entry.Price) {
		return
	}

	player.BuyHull(entry.Hull, entry.Price)
}

func (us *UpgradeSystem) tryUpgradeFuelTank(player *entities.Player) {
	entry := us.fuelTankShop.GetNextFuelTank(player.FuelTank.Tier())
	if entry == nil {
		return // Max level reached
	}

	if !player.CanAfford(entry.Price) {
		return
	}

	player.BuyFuelTank(entry.FuelTank, entry.Price)
}

func (us *UpgradeSystem) tryUpgradeCargoHold(player *entities.Player) {
	entry := us.cargoHoldShop.GetNextCargoHold(player.CargoHold.Tier())
	if entry == nil {
		return // Max level reached
	}

	if !player.CanAfford(entry.Price) {
		return
	}

	player.BuyCargoHold(entry.CargoHold, entry.Price)
}

func (us *UpgradeSystem) tryUpgradeHeatShield(player *entities.Player) {
	entry := us.heatShieldShop.GetNextHeatShield(player.HeatShield.Tier())
	if entry == nil {
		return // Max level reached
	}

	if !player.CanAfford(entry.Price) {
		return
	}

	player.BuyHeatShield(entry.HeatShield, entry.Price)
}

func (us *UpgradeSystem) GetEngineShop() *entities.EngineUpgradeShop {
	return us.engineShop
}

func (us *UpgradeSystem) GetHullShop() *entities.HullUpgradeShop {
	return us.hullShop
}

func (us *UpgradeSystem) GetFuelTankShop() *entities.FuelTankUpgradeShop {
	return us.fuelTankShop
}

func (us *UpgradeSystem) GetCargoHoldShop() *entities.CargoHoldUpgradeShop {
	return us.cargoHoldShop
}

func (us *UpgradeSystem) GetHeatShieldShop() *entities.HeatShieldUpgradeShop {
	return us.heatShieldShop
}
