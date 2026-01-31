package phases

// Config defines a boss phase
type Config struct {
	HPThreshold        float32 // Percentage of HP where this phase ends (e.g., 0.66 = 66%)
	MovementSpeed      float32 // Movement speed in this phase
	ProjectileCooldown float32 // Time between projectile attacks
	AOECooldown        float32 // Time between AOE attacks (0 = disabled)
}

// Manager tracks boss phases based on HP
type Manager struct {
	phases       []Config
	currentPhase int
	maxHP        float32
}

// NewManager creates a phase manager with the given phases
// Phases should be ordered from full HP to low HP
func NewManager(maxHP float32, phases []Config) *Manager {
	return &Manager{
		phases:       phases,
		currentPhase: 0,
		maxHP:        maxHP,
	}
}

// Update checks if phase should change based on current HP
// Returns true if phase changed
func (pm *Manager) Update(currentHP float32) bool {
	if len(pm.phases) == 0 {
		return false
	}

	hpPercent := currentHP / pm.maxHP
	oldPhase := pm.currentPhase

	// Check each phase threshold
	for i := pm.currentPhase; i < len(pm.phases)-1; i++ {
		if hpPercent <= pm.phases[i].HPThreshold {
			pm.currentPhase = i + 1
		}
	}

	return pm.currentPhase != oldPhase
}

func (pm *Manager) GetCurrentPhase() int {
	return pm.currentPhase
}

func (pm *Manager) GetCurrentConfig() Config {
	if pm.currentPhase >= len(pm.phases) {
		return pm.phases[len(pm.phases)-1]
	}
	return pm.phases[pm.currentPhase]
}
