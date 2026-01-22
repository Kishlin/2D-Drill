package entities

import "github.com/Kishlin/drill-game/internal/domain/components"

const (
	BuildingWidth  = 320.0
	BuildingHeight = 192.0
)

type Building struct {
	Position     components.Position
	Interactable components.Interactable
}

func NewBuilding(x, y, width, height float32, interactableType components.InteractableType) *Building {
	return &Building{
		Position:     components.NewPosition(x, y, width, height),
		Interactable: components.Interactable{Type: interactableType},
	}
}

func NewMarketBuilding(x, y float32) *Building {
	return NewBuilding(x, y, BuildingWidth, BuildingHeight, components.InteractableMarket)
}

func NewFuelStationBuilding(x, y float32) *Building {
	return NewBuilding(x, y, BuildingWidth, BuildingHeight, components.InteractableFuelStation)
}

func NewHospitalBuilding(x, y float32) *Building {
	return NewBuilding(x, y, BuildingWidth, BuildingHeight, components.InteractableHospital)
}

func NewUpgradeShopBuilding(x, y float32) *Building {
	return NewBuilding(x, y, BuildingWidth, BuildingHeight, components.InteractableUpgradeShop)
}

func NewItemShopBuilding(x, y float32) *Building {
	return NewBuilding(x, y, BuildingWidth, BuildingHeight, components.InteractableItemShop)
}
