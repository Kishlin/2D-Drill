package effects

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/entities"
)

func TestClearOreInventory_Apply(t *testing.T) {
	player := testPlayer()
	player.OreInventory = map[string]int{
		"copper": 5,
		"gold":   3,
		"iron":   10,
	}

	effect := ClearOreInventory{}
	effect.Apply(testContext(player))

	if len(player.OreInventory) != 0 {
		t.Errorf("expected ore inventory to be empty, got %d items", len(player.OreInventory))
	}
}

func TestClearOreInventory_Apply_AlreadyEmpty(t *testing.T) {
	player := testPlayer()
	player.OreInventory = make(map[string]int)

	effect := ClearOreInventory{}
	effect.Apply(testContext(player))

	if len(player.OreInventory) != 0 {
		t.Errorf("expected ore inventory to remain empty, got %d items", len(player.OreInventory))
	}
}

func TestClearOreInventory_Apply_CreatesNewMap(t *testing.T) {
	player := testPlayer()
	oldInventory := player.OreInventory
	oldInventory["copper"] = 5

	effect := ClearOreInventory{}
	effect.Apply(testContext(player))

	// Add to old inventory shouldn't affect player
	oldInventory["gold"] = 10

	if len(player.OreInventory) != 0 {
		t.Errorf("expected ore inventory to be independent of old reference")
	}
}

func TestAddItem_Apply_Teleport(t *testing.T) {
	player := testPlayer()
	player.ItemInventory = [5]int{0, 0, 0, 0, 0}

	effect := AddItem{ItemType: entities.ItemTeleport}
	effect.Apply(testContext(player))

	if player.ItemInventory[entities.ItemTeleport] != 1 {
		t.Errorf("expected teleport count to be 1, got %d", player.ItemInventory[entities.ItemTeleport])
	}
}

func TestAddItem_Apply_Repair(t *testing.T) {
	player := testPlayer()
	player.ItemInventory = [5]int{0, 0, 0, 0, 0}

	effect := AddItem{ItemType: entities.ItemRepair}
	effect.Apply(testContext(player))

	if player.ItemInventory[entities.ItemRepair] != 1 {
		t.Errorf("expected repair count to be 1, got %d", player.ItemInventory[entities.ItemRepair])
	}
}

func TestAddItem_Apply_Refuel(t *testing.T) {
	player := testPlayer()
	player.ItemInventory = [5]int{0, 0, 0, 0, 0}

	effect := AddItem{ItemType: entities.ItemRefuel}
	effect.Apply(testContext(player))

	if player.ItemInventory[entities.ItemRefuel] != 1 {
		t.Errorf("expected refuel count to be 1, got %d", player.ItemInventory[entities.ItemRefuel])
	}
}

func TestAddItem_Apply_Bomb(t *testing.T) {
	player := testPlayer()
	player.ItemInventory = [5]int{0, 0, 0, 0, 0}

	effect := AddItem{ItemType: entities.ItemBomb}
	effect.Apply(testContext(player))

	if player.ItemInventory[entities.ItemBomb] != 1 {
		t.Errorf("expected bomb count to be 1, got %d", player.ItemInventory[entities.ItemBomb])
	}
}

func TestAddItem_Apply_BigBomb(t *testing.T) {
	player := testPlayer()
	player.ItemInventory = [5]int{0, 0, 0, 0, 0}

	effect := AddItem{ItemType: entities.ItemBigBomb}
	effect.Apply(testContext(player))

	if player.ItemInventory[entities.ItemBigBomb] != 1 {
		t.Errorf("expected big bomb count to be 1, got %d", player.ItemInventory[entities.ItemBigBomb])
	}
}

func TestAddItem_Apply_StacksWithExisting(t *testing.T) {
	player := testPlayer()
	player.ItemInventory = [5]int{3, 0, 0, 0, 0} // 3 teleports already

	effect := AddItem{ItemType: entities.ItemTeleport}
	effect.Apply(testContext(player))

	if player.ItemInventory[entities.ItemTeleport] != 4 {
		t.Errorf("expected teleport count to be 4, got %d", player.ItemInventory[entities.ItemTeleport])
	}
}

func TestAddItem_Apply_MultipleItems(t *testing.T) {
	player := testPlayer()
	ctx := testContext(player)
	player.ItemInventory = [5]int{0, 0, 0, 0, 0}

	AddItem{ItemType: entities.ItemTeleport}.Apply(ctx)
	AddItem{ItemType: entities.ItemBomb}.Apply(ctx)
	AddItem{ItemType: entities.ItemTeleport}.Apply(ctx)

	if player.ItemInventory[entities.ItemTeleport] != 2 {
		t.Errorf("expected teleport count to be 2, got %d", player.ItemInventory[entities.ItemTeleport])
	}
	if player.ItemInventory[entities.ItemBomb] != 1 {
		t.Errorf("expected bomb count to be 1, got %d", player.ItemInventory[entities.ItemBomb])
	}
}
