package sentinel_boss

import "github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"

// State IDs
const (
	StateHover statemachine.StateID = iota
	StateChargeWindup
	StateCharge
	StateStunned
	StateLaserAim
	StateLaser
)
