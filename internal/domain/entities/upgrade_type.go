package entities

type UpgradeType int

const (
	UpgradeEngine UpgradeType = iota
	UpgradeHull
	UpgradeFuelTank
	UpgradeCargoHold
	UpgradeHeatShield
	UpgradeDrill
	UpgradeTypeCount // Total number of upgrade types (6)
)

func (ut UpgradeType) String() string {
	switch ut {
	case UpgradeEngine:
		return "Engine"
	case UpgradeHull:
		return "Hull"
	case UpgradeFuelTank:
		return "Fuel Tank"
	case UpgradeCargoHold:
		return "Cargo Hold"
	case UpgradeHeatShield:
		return "Heat Shield"
	case UpgradeDrill:
		return "Drill"
	default:
		return "Unknown"
	}
}

func (ut UpgradeType) ShortName() string {
	switch ut {
	case UpgradeEngine:
		return "Engine"
	case UpgradeHull:
		return "Hull"
	case UpgradeFuelTank:
		return "Fuel"
	case UpgradeCargoHold:
		return "Cargo"
	case UpgradeHeatShield:
		return "Heat"
	case UpgradeDrill:
		return "Drill"
	default:
		return "?"
	}
}
