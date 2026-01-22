package components

type InteractableType int

const (
	InteractableMarket InteractableType = iota
	InteractableFuelStation
	InteractableHospital
	InteractableUpgradeShop
	InteractableItemShop
)

type Interactable struct {
	Type InteractableType
}
