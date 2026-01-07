package entities

import "testing"

func TestPlayer_AddOre_SingleType(t *testing.T) {
	player := NewPlayer(0, 0)

	success := player.AddOre(OreCopper)

	if !success {
		t.Errorf("Expected AddOre to succeed")
	}
	if player.OreInventory[OreCopper] != 1 {
		t.Errorf("Expected 1 copper, got %d", player.OreInventory[OreCopper])
	}
}

func TestPlayer_AddOre_MultipleTypes(t *testing.T) {
	player := NewPlayer(0, 0)

	player.AddOre(OreCopper)
	player.AddOre(OreCopper)
	player.AddOre(OreCopper)
	player.AddOre(OreGold)
	player.AddOre(OreGold)
	player.AddOre(OreGold)
	player.AddOre(OreGold)
	player.AddOre(OreDiamond)

	if player.OreInventory[OreCopper] != 3 {
		t.Errorf("Expected 3 copper, got %d", player.OreInventory[OreCopper])
	}
	if player.OreInventory[OreGold] != 4 {
		t.Errorf("Expected 4 gold, got %d", player.OreInventory[OreGold])
	}
	if player.OreInventory[OreDiamond] != 1 {
		t.Errorf("Expected 1 diamond, got %d", player.OreInventory[OreDiamond])
	}
}

func TestPlayer_AddOre_Accumulates(t *testing.T) {
	player := NewPlayer(0, 0)

	for i := 0; i < 10; i++ {
		player.AddOre(OreIron)
	}

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

	// Should return false on invalid ore types and not panic
	if player.AddOre(OreType(-1)) {
		t.Errorf("Should return false for invalid ore type")
	}
	if player.AddOre(OreType(999)) {
		t.Errorf("Should return false for invalid ore type")
	}

	// Verify inventory is still all zeros
	for _, oreType := range GetAllOreTypes() {
		if player.OreInventory[oreType] != 0 {
			t.Errorf("Invalid ore types should not affect inventory")
		}
	}
}

func TestPlayer_AddOre_CargoCapacity(t *testing.T) {
	player := NewPlayer(0, 0)
	// Player starts with Base CargoHold, capacity 10

	// Fill cargo to capacity
	for i := 0; i < 10; i++ {
		if !player.AddOre(OreCopper) {
			t.Errorf("AddOre should succeed at position %d", i)
		}
	}

	if player.GetTotalOreCount() != 10 {
		t.Errorf("Expected 10 total ore, got %d", player.GetTotalOreCount())
	}

	// Next ore should fail
	if player.AddOre(OreCopper) {
		t.Errorf("AddOre should fail when cargo is full")
	}

	if player.GetTotalOreCount() != 10 {
		t.Errorf("Failed AddOre should not change inventory")
	}
}

// DealDamage tests

func TestPlayer_DealDamage_ReducesHP(t *testing.T) {
	player := NewPlayer(0, 0)
	initialHP := player.HP

	player.DealDamage(2.0)

	if player.HP != initialHP-2.0 {
		t.Errorf("Expected HP %f, got %f", initialHP-2.0, player.HP)
	}
}

func TestPlayer_DealDamage_SmallDamage(t *testing.T) {
	player := NewPlayer(0, 0)
	// Player starts with 10 HP

	player.DealDamage(1.5)

	if player.HP != 8.5 {
		t.Errorf("Expected 8.5 HP, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_LethalDamage(t *testing.T) {
	player := NewPlayer(0, 0)

	player.DealDamage(10.0)

	if player.HP != 0.0 {
		t.Errorf("Expected 0 HP, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_OverDamage(t *testing.T) {
	player := NewPlayer(0, 0)

	// Deal more damage than current HP
	player.DealDamage(100.0)

	// Should clamp at 0, not go negative
	if player.HP != 0.0 {
		t.Errorf("Expected HP clamped at 0, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_MultipleDamageInstances(t *testing.T) {
	player := NewPlayer(0, 0)
	// Player starts with 10 HP

	player.DealDamage(2.0)
	player.DealDamage(3.0)
	player.DealDamage(1.0)

	if player.HP != 4.0 {
		t.Errorf("Expected 4 HP after 3 damage instances, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_AlreadyDead(t *testing.T) {
	player := NewPlayer(0, 0)

	// Kill player
	player.DealDamage(10.0)

	// Deal additional damage
	player.DealDamage(5.0)

	// Should remain at 0, not go negative
	if player.HP != 0.0 {
		t.Errorf("Expected HP to stay at 0 for dead player, got %f", player.HP)
	}
}

func TestPlayer_DealDamage_PartialDamage(t *testing.T) {
	player := NewPlayer(0, 0)
	// Player starts with 10 HP

	player.DealDamage(3.7)

	expectedHP := float32(10.0 - 3.7)
	if player.HP != expectedHP {
		t.Errorf("Expected %f HP, got %f", expectedHP, player.HP)
	}
}

func TestPlayer_DealDamage_ZeroDamage(t *testing.T) {
	player := NewPlayer(0, 0)
	initialHP := player.HP

	player.DealDamage(0.0)

	if player.HP != initialHP {
		t.Errorf("Expected no change with zero damage, got HP %f", player.HP)
	}
}
