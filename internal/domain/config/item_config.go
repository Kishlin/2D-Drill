package config

type ItemEntry struct {
	Name   string
	Price  int
	Radius int // For bombs only (0 for non-bombs)
}

type ItemConfig struct {
	Teleport ItemEntry
	Repair   ItemEntry
	Refuel   ItemEntry
	Bomb     ItemEntry
	BigBomb  ItemEntry
}

func DefaultItemConfig() ItemConfig {
	return ItemConfig{
		Teleport: ItemEntry{Name: "Teleport", Price: 500, Radius: 0},
		Repair:   ItemEntry{Name: "Repair Kit", Price: 200, Radius: 0},
		Refuel:   ItemEntry{Name: "Fuel Can", Price: 100, Radius: 0},
		Bomb:     ItemEntry{Name: "Bomb", Price: 300, Radius: 1},     // 3x3 area
		BigBomb:  ItemEntry{Name: "Big Bomb", Price: 800, Radius: 2}, // 5x5 area
	}
}

func (c *ItemConfig) GetItemByIndex(index int) *ItemEntry {
	switch index {
	case 0:
		return &c.Teleport
	case 1:
		return &c.Repair
	case 2:
		return &c.Refuel
	case 3:
		return &c.Bomb
	case 4:
		return &c.BigBomb
	default:
		return nil
	}
}

func (c *ItemConfig) GetAllItems() []ItemEntry {
	return []ItemEntry{c.Teleport, c.Repair, c.Refuel, c.Bomb, c.BigBomb}
}
