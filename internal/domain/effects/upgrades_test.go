package effects

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/upgrades"
)

func TestSetUpgrade_Engine(t *testing.T) {
	player := testPlayer()

	newEngine := upgrades.NewEngineFromConfig(1, "Mk1 Engine", config.EngineStats{
		MaxSpeed: 150, Acceleration: 75, FlyAcceleration: 45, MaxUpwardSpeed: -300,
	})

	effect := SetUpgrade{Upgrade: newEngine}
	effect.Apply(testContext(player))

	if player.GetUpgradeTier(upgrades.TypeEngine) != 1 {
		t.Errorf("expected engine tier to be 1, got %d", player.GetUpgradeTier(upgrades.TypeEngine))
	}
	if player.GetUpgrade(upgrades.TypeEngine).Name() != "Mk1 Engine" {
		t.Errorf("expected engine name to be 'Mk1 Engine', got '%s'", player.GetUpgrade(upgrades.TypeEngine).Name())
	}
	if player.MaxSpeed() != 150 {
		t.Errorf("expected max speed to be 150, got %.2f", player.MaxSpeed())
	}
}

func TestSetUpgrade_Hull(t *testing.T) {
	player := testPlayer()

	newHull := upgrades.NewHullFromConfig(2, "Mk2 Hull", config.HullStats{MaxHP: 200})

	effect := SetUpgrade{Upgrade: newHull}
	effect.Apply(testContext(player))

	if player.GetUpgradeTier(upgrades.TypeHull) != 2 {
		t.Errorf("expected hull tier to be 2, got %d", player.GetUpgradeTier(upgrades.TypeHull))
	}
	if player.GetUpgrade(upgrades.TypeHull).Name() != "Mk2 Hull" {
		t.Errorf("expected hull name to be 'Mk2 Hull', got '%s'", player.GetUpgrade(upgrades.TypeHull).Name())
	}
	if player.MaxHP() != 200 {
		t.Errorf("expected max HP to be 200, got %.2f", player.MaxHP())
	}
}

func TestSetUpgrade_FuelTank(t *testing.T) {
	player := testPlayer()

	newTank := upgrades.NewFuelTankFromConfig(3, "Mk3 Tank", config.FuelTankStats{Capacity: 250})

	effect := SetUpgrade{Upgrade: newTank}
	effect.Apply(testContext(player))

	if player.GetUpgradeTier(upgrades.TypeFuelTank) != 3 {
		t.Errorf("expected fuel tank tier to be 3, got %d", player.GetUpgradeTier(upgrades.TypeFuelTank))
	}
	if player.GetUpgrade(upgrades.TypeFuelTank).Name() != "Mk3 Tank" {
		t.Errorf("expected fuel tank name to be 'Mk3 Tank', got '%s'", player.GetUpgrade(upgrades.TypeFuelTank).Name())
	}
	if player.FuelCapacity() != 250 {
		t.Errorf("expected capacity to be 250, got %.2f", player.FuelCapacity())
	}
}

func TestSetUpgrade_CargoHold(t *testing.T) {
	player := testPlayer()

	newCargo := upgrades.NewCargoHoldFromConfig(4, "Mk4 Cargo", config.CargoHoldStats{Capacity: 50})

	effect := SetUpgrade{Upgrade: newCargo}
	effect.Apply(testContext(player))

	if player.GetUpgradeTier(upgrades.TypeCargoHold) != 4 {
		t.Errorf("expected cargo hold tier to be 4, got %d", player.GetUpgradeTier(upgrades.TypeCargoHold))
	}
	if player.GetUpgrade(upgrades.TypeCargoHold).Name() != "Mk4 Cargo" {
		t.Errorf("expected cargo hold name to be 'Mk4 Cargo', got '%s'", player.GetUpgrade(upgrades.TypeCargoHold).Name())
	}
	if player.CargoCapacity() != 50 {
		t.Errorf("expected capacity to be 50, got %d", player.CargoCapacity())
	}
}

func TestSetUpgrade_HeatShield(t *testing.T) {
	player := testPlayer()

	newShield := upgrades.NewHeatShieldFromConfig(5, "Mk5 Shield", config.HeatShieldStats{HeatResistance: 500})

	effect := SetUpgrade{Upgrade: newShield}
	effect.Apply(testContext(player))

	if player.GetUpgradeTier(upgrades.TypeHeatShield) != 5 {
		t.Errorf("expected heat shield tier to be 5, got %d", player.GetUpgradeTier(upgrades.TypeHeatShield))
	}
	if player.GetUpgrade(upgrades.TypeHeatShield).Name() != "Mk5 Shield" {
		t.Errorf("expected heat shield name to be 'Mk5 Shield', got '%s'", player.GetUpgrade(upgrades.TypeHeatShield).Name())
	}
	if player.HeatResistance() != 500 {
		t.Errorf("expected heat resistance to be 500, got %.2f", player.HeatResistance())
	}
}

func TestSetUpgrade_Drill(t *testing.T) {
	player := testPlayer()

	newDrill := upgrades.NewDrillFromConfig(2, "Mk2 Drill", config.DrillStats{DrillSpeed: 2.5})

	effect := SetUpgrade{Upgrade: newDrill}
	effect.Apply(testContext(player))

	if player.GetUpgradeTier(upgrades.TypeDrill) != 2 {
		t.Errorf("expected drill tier to be 2, got %d", player.GetUpgradeTier(upgrades.TypeDrill))
	}
	if player.GetUpgrade(upgrades.TypeDrill).Name() != "Mk2 Drill" {
		t.Errorf("expected drill name to be 'Mk2 Drill', got '%s'", player.GetUpgrade(upgrades.TypeDrill).Name())
	}
	if player.DrillSpeed() != 2.5 {
		t.Errorf("expected drill speed to be 2.5, got %.2f", player.DrillSpeed())
	}
}
