package effects

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
	"github.com/Kishlin/drill-game/internal/domain/entities"
)

func TestSetEngine_Apply(t *testing.T) {
	player := testPlayer()

	oldEngine := entities.NewEngineFromConfig(0, "Base Engine", config.EngineStats{
		MaxSpeed: 100, Acceleration: 50, FlyAcceleration: 30, MaxUpwardSpeed: -200,
	})
	newEngine := entities.NewEngineFromConfig(1, "Mk1 Engine", config.EngineStats{
		MaxSpeed: 150, Acceleration: 75, FlyAcceleration: 45, MaxUpwardSpeed: -300,
	})

	player.Engine = oldEngine

	effect := SetEngine{Engine: newEngine}
	effect.Apply(player)

	if player.Engine.Tier() != 1 {
		t.Errorf("expected engine tier to be 1, got %d", player.Engine.Tier())
	}
	if player.Engine.Name() != "Mk1 Engine" {
		t.Errorf("expected engine name to be 'Mk1 Engine', got '%s'", player.Engine.Name())
	}
	if player.Engine.MaxSpeed() != 150 {
		t.Errorf("expected max speed to be 150, got %.2f", player.Engine.MaxSpeed())
	}
}

func TestSetHull_Apply(t *testing.T) {
	player := testPlayer()

	oldHull := entities.NewHullFromConfig(0, "Base Hull", config.HullStats{MaxHP: 100})
	newHull := entities.NewHullFromConfig(2, "Mk2 Hull", config.HullStats{MaxHP: 200})

	player.Hull = oldHull

	effect := SetHull{Hull: newHull}
	effect.Apply(player)

	if player.Hull.Tier() != 2 {
		t.Errorf("expected hull tier to be 2, got %d", player.Hull.Tier())
	}
	if player.Hull.Name() != "Mk2 Hull" {
		t.Errorf("expected hull name to be 'Mk2 Hull', got '%s'", player.Hull.Name())
	}
	if player.Hull.MaxHP() != 200 {
		t.Errorf("expected max HP to be 200, got %.2f", player.Hull.MaxHP())
	}
}

func TestSetFuelTank_Apply(t *testing.T) {
	player := testPlayer()

	oldTank := entities.NewFuelTankFromConfig(0, "Base Tank", config.FuelTankStats{Capacity: 100})
	newTank := entities.NewFuelTankFromConfig(3, "Mk3 Tank", config.FuelTankStats{Capacity: 250})

	player.FuelTank = oldTank

	effect := SetFuelTank{FuelTank: newTank}
	effect.Apply(player)

	if player.FuelTank.Tier() != 3 {
		t.Errorf("expected fuel tank tier to be 3, got %d", player.FuelTank.Tier())
	}
	if player.FuelTank.Name() != "Mk3 Tank" {
		t.Errorf("expected fuel tank name to be 'Mk3 Tank', got '%s'", player.FuelTank.Name())
	}
	if player.FuelTank.Capacity() != 250 {
		t.Errorf("expected capacity to be 250, got %.2f", player.FuelTank.Capacity())
	}
}

func TestSetCargoHold_Apply(t *testing.T) {
	player := testPlayer()

	oldCargo := entities.NewCargoHoldFromConfig(0, "Base Cargo", config.CargoHoldStats{Capacity: 10})
	newCargo := entities.NewCargoHoldFromConfig(4, "Mk4 Cargo", config.CargoHoldStats{Capacity: 50})

	player.CargoHold = oldCargo

	effect := SetCargoHold{CargoHold: newCargo}
	effect.Apply(player)

	if player.CargoHold.Tier() != 4 {
		t.Errorf("expected cargo hold tier to be 4, got %d", player.CargoHold.Tier())
	}
	if player.CargoHold.Name() != "Mk4 Cargo" {
		t.Errorf("expected cargo hold name to be 'Mk4 Cargo', got '%s'", player.CargoHold.Name())
	}
	if player.CargoHold.Capacity() != 50 {
		t.Errorf("expected capacity to be 50, got %d", player.CargoHold.Capacity())
	}
}

func TestSetHeatShield_Apply(t *testing.T) {
	player := testPlayer()

	oldShield := entities.NewHeatShieldFromConfig(0, "Base Shield", config.HeatShieldStats{HeatResistance: 100})
	newShield := entities.NewHeatShieldFromConfig(5, "Mk5 Shield", config.HeatShieldStats{HeatResistance: 500})

	player.HeatShield = oldShield

	effect := SetHeatShield{HeatShield: newShield}
	effect.Apply(player)

	if player.HeatShield.Tier() != 5 {
		t.Errorf("expected heat shield tier to be 5, got %d", player.HeatShield.Tier())
	}
	if player.HeatShield.Name() != "Mk5 Shield" {
		t.Errorf("expected heat shield name to be 'Mk5 Shield', got '%s'", player.HeatShield.Name())
	}
	if player.HeatShield.HeatResistance() != 500 {
		t.Errorf("expected heat resistance to be 500, got %.2f", player.HeatShield.HeatResistance())
	}
}

func TestSetDrill_Apply(t *testing.T) {
	player := testPlayer()

	oldDrill := entities.NewDrillFromConfig(0, "Base Drill", config.DrillStats{DrillSpeed: 1.0})
	newDrill := entities.NewDrillFromConfig(2, "Mk2 Drill", config.DrillStats{DrillSpeed: 2.5})

	player.Drill = oldDrill

	effect := SetDrill{Drill: newDrill}
	effect.Apply(player)

	if player.Drill.Tier() != 2 {
		t.Errorf("expected drill tier to be 2, got %d", player.Drill.Tier())
	}
	if player.Drill.Name() != "Mk2 Drill" {
		t.Errorf("expected drill name to be 'Mk2 Drill', got '%s'", player.Drill.Name())
	}
	if player.Drill.DrillSpeed() != 2.5 {
		t.Errorf("expected drill speed to be 2.5, got %.2f", player.Drill.DrillSpeed())
	}
}
