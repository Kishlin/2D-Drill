package effects

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

func testPlayer() *entities.Player {
	return &entities.Player{
		AABB:         types.NewAABB(0, 0, 54, 54),
		Money:        100,
		Fuel:         50.0,
		HP:           75.0,
		OreInventory: make(map[string]int),
	}
}

func TestTakeMoney_Apply(t *testing.T) {
	player := testPlayer()
	player.Money = 100

	effect := TakeMoney{Amount: 30}
	effect.Apply(player)

	if player.Money != 70 {
		t.Errorf("expected money to be 70, got %d", player.Money)
	}
}

func TestTakeMoney_Apply_TakesAllMoney(t *testing.T) {
	player := testPlayer()
	player.Money = 50

	effect := TakeMoney{Amount: 50}
	effect.Apply(player)

	if player.Money != 0 {
		t.Errorf("expected money to be 0, got %d", player.Money)
	}
}

func TestTakeMoney_Apply_CanGoNegative(t *testing.T) {
	player := testPlayer()
	player.Money = 20

	effect := TakeMoney{Amount: 50}
	effect.Apply(player)

	if player.Money != -30 {
		t.Errorf("expected money to be -30, got %d", player.Money)
	}
}

func TestAddMoney_Apply(t *testing.T) {
	player := testPlayer()
	player.Money = 100

	effect := AddMoney{Amount: 50}
	effect.Apply(player)

	if player.Money != 150 {
		t.Errorf("expected money to be 150, got %d", player.Money)
	}
}

func TestAddMoney_Apply_FromZero(t *testing.T) {
	player := testPlayer()
	player.Money = 0

	effect := AddMoney{Amount: 25}
	effect.Apply(player)

	if player.Money != 25 {
		t.Errorf("expected money to be 25, got %d", player.Money)
	}
}

func TestAddMoney_Apply_LargeAmount(t *testing.T) {
	player := testPlayer()
	player.Money = 100

	effect := AddMoney{Amount: 10000}
	effect.Apply(player)

	if player.Money != 10100 {
		t.Errorf("expected money to be 10100, got %d", player.Money)
	}
}
