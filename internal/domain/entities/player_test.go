package entities

import "testing"

func TestPlayer_AddOre_SingleType(t *testing.T) {
	player := NewPlayer(0, 0)

	player.AddOre(OreCopper, 5)

	if player.OreInventory[OreCopper] != 5 {
		t.Errorf("Expected 5 copper, got %d", player.OreInventory[OreCopper])
	}
}

func TestPlayer_AddOre_MultipleTypes(t *testing.T) {
	player := NewPlayer(0, 0)

	player.AddOre(OreCopper, 3)
	player.AddOre(OreGold, 7)
	player.AddOre(OreDiamond, 1)

	if player.OreInventory[OreCopper] != 3 {
		t.Errorf("Expected 3 copper, got %d", player.OreInventory[OreCopper])
	}
	if player.OreInventory[OreGold] != 7 {
		t.Errorf("Expected 7 gold, got %d", player.OreInventory[OreGold])
	}
	if player.OreInventory[OreDiamond] != 1 {
		t.Errorf("Expected 1 diamond, got %d", player.OreInventory[OreDiamond])
	}
}

func TestPlayer_AddOre_Accumulates(t *testing.T) {
	player := NewPlayer(0, 0)

	player.AddOre(OreIron, 5)
	player.AddOre(OreIron, 3)
	player.AddOre(OreIron, 2)

	if player.OreInventory[OreIron] != 10 {
		t.Errorf("Expected 10 iron, got %d", player.OreInventory[OreIron])
	}
}

func TestPlayer_NewPlayer_StartsWithZeroOres(t *testing.T) {
	player := NewPlayer(0, 0)

	for _, oreType := range GetAllOreTypes() {
		if player.OreInventory[oreType] != 0 {
			t.Errorf("New player should have 0 of ore type %d, got %d", oreType, player.OreInventory[oreType])
		}
	}
}

func TestPlayer_AddOre_BoundsCheck(t *testing.T) {
	player := NewPlayer(0, 0)

	// Should not panic on invalid ore types
	player.AddOre(OreType(-1), 5)
	player.AddOre(OreType(999), 5)

	// Verify inventory is still all zeros
	for _, oreType := range GetAllOreTypes() {
		if player.OreInventory[oreType] != 0 {
			t.Errorf("Invalid ore types should not affect inventory")
		}
	}
}
