package entities

type ItemType int

const (
	ItemTeleport ItemType = iota
	ItemRepair
	ItemRefuel
	ItemBomb
	ItemBigBomb
)

// ItemNames provides display names for each item type
var ItemNames = map[ItemType]string{
	ItemTeleport: "Teleport",
	ItemRepair:   "Repair Kit",
	ItemRefuel:   "Fuel Can",
	ItemBomb:     "Bomb",
	ItemBigBomb:  "Big Bomb",
}
