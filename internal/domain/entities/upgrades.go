package entities

type UpgradeLevel int

//goland:noinspection GoUnusedConst
const (
	UpgradeLevelBase UpgradeLevel = 0
	UpgradeLevelMk1  UpgradeLevel = 1
	UpgradeLevelMk2  UpgradeLevel = 2
	UpgradeLevelMk3  UpgradeLevel = 3
	UpgradeLevelMk4  UpgradeLevel = 4
	UpgradeLevelMk5  UpgradeLevel = 5 // Max level
)

type Upgrades struct {
	Engine   UpgradeLevel
	Hull     UpgradeLevel
	FuelTank UpgradeLevel
}

func NewUpgrades() Upgrades {
	return Upgrades{
		Engine:   UpgradeLevelBase,
		Hull:     UpgradeLevelBase,
		FuelTank: UpgradeLevelBase,
	}
}
