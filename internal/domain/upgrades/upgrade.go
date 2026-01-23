package upgrades

type UpgradeType int

const (
	TypeEngine UpgradeType = iota
	TypeHull
	TypeFuelTank
	TypeCargoHold
	TypeHeatShield
	TypeDrill
	TypeCount
)

var typeNames = [TypeCount]string{"Engine", "Hull", "Fuel Tank", "Cargo Hold", "Heat Shield", "Drill"}
var typeShortNames = [TypeCount]string{"Engine", "Hull", "Fuel", "Cargo", "Heat", "Drill"}

func (t UpgradeType) String() string {
	if t < 0 || t >= TypeCount {
		return "Unknown"
	}
	return typeNames[t]
}

func (t UpgradeType) ShortName() string {
	if t < 0 || t >= TypeCount {
		return "?"
	}
	return typeShortNames[t]
}

type Upgrade interface {
	Tier() int
	Name() string
	Type() UpgradeType
}
