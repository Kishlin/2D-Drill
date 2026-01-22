package effects

import (
	"testing"
)

func TestSetFuel_Apply(t *testing.T) {
	player := testPlayer()
	player.Fuel = 25.0

	effect := SetFuel{Amount: 100.0}
	effect.Apply(player)

	if player.Fuel != 100.0 {
		t.Errorf("expected fuel to be 100.0, got %.2f", player.Fuel)
	}
}

func TestSetFuel_Apply_ToZero(t *testing.T) {
	player := testPlayer()
	player.Fuel = 50.0

	effect := SetFuel{Amount: 0.0}
	effect.Apply(player)

	if player.Fuel != 0.0 {
		t.Errorf("expected fuel to be 0.0, got %.2f", player.Fuel)
	}
}

func TestSetFuel_Apply_PartialRefuel(t *testing.T) {
	player := testPlayer()
	player.Fuel = 10.0

	effect := SetFuel{Amount: 75.5}
	effect.Apply(player)

	if player.Fuel != 75.5 {
		t.Errorf("expected fuel to be 75.5, got %.2f", player.Fuel)
	}
}

func TestSetHP_Apply(t *testing.T) {
	player := testPlayer()
	player.HP = 50.0

	effect := SetHP{Amount: 100.0}
	effect.Apply(player)

	if player.HP != 100.0 {
		t.Errorf("expected HP to be 100.0, got %.2f", player.HP)
	}
}

func TestSetHP_Apply_ToZero(t *testing.T) {
	player := testPlayer()
	player.HP = 75.0

	effect := SetHP{Amount: 0.0}
	effect.Apply(player)

	if player.HP != 0.0 {
		t.Errorf("expected HP to be 0.0, got %.2f", player.HP)
	}
}

func TestSetHP_Apply_PartialHeal(t *testing.T) {
	player := testPlayer()
	player.HP = 25.0

	effect := SetHP{Amount: 80.0}
	effect.Apply(player)

	if player.HP != 80.0 {
		t.Errorf("expected HP to be 80.0, got %.2f", player.HP)
	}
}
