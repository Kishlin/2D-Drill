package config

type ItemEntry struct {
	Name   string
	Price  int
	Radius int     // For bombs only (0 for non-bombs)
	Damage float32 // For bombs only (0 for non-bombs)
}

type ItemConfig struct {
	Teleport ItemEntry
	Repair   ItemEntry
	Refuel   ItemEntry
	Bomb     ItemEntry
	BigBomb  ItemEntry
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
