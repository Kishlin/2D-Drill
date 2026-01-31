package phases

import "testing"

func TestManager_StartsAtPhase0(t *testing.T) {
	p := []Config{
		{HPThreshold: 0.66},
		{HPThreshold: 0.33},
		{HPThreshold: 0.0},
	}
	pm := NewManager(100.0, p)

	if pm.GetCurrentPhase() != 0 {
		t.Errorf("expected phase 0, got %d", pm.GetCurrentPhase())
	}
}

func TestManager_TransitionsToPhase1(t *testing.T) {
	p := []Config{
		{HPThreshold: 0.66},
		{HPThreshold: 0.33},
		{HPThreshold: 0.0},
	}
	pm := NewManager(100.0, p)

	// HP drops below 66%
	changed := pm.Update(65.0)

	if changed == false {
		t.Error("expected phase change")
	}

	if pm.GetCurrentPhase() != 1 {
		t.Errorf("expected phase 1, got %d", pm.GetCurrentPhase())
	}
}

func TestManager_TransitionsToPhase2(t *testing.T) {
	p := []Config{
		{HPThreshold: 0.66},
		{HPThreshold: 0.33},
		{HPThreshold: 0.0},
	}
	pm := NewManager(100.0, p)

	// HP drops below 33%
	pm.Update(32.0)

	if pm.GetCurrentPhase() != 2 {
		t.Errorf("expected phase 2, got %d", pm.GetCurrentPhase())
	}
}

func TestManager_NoChangeAboveThreshold(t *testing.T) {
	p := []Config{
		{HPThreshold: 0.66},
		{HPThreshold: 0.33},
	}
	pm := NewManager(100.0, p)

	changed := pm.Update(80.0)

	if changed {
		t.Error("expected no phase change")
	}

	if pm.GetCurrentPhase() != 0 {
		t.Errorf("expected phase 0, got %d", pm.GetCurrentPhase())
	}
}

func TestManager_GetCurrentConfig(t *testing.T) {
	p := []Config{
		{HPThreshold: 0.66, MovementSpeed: 80.0},
		{HPThreshold: 0.33, MovementSpeed: 100.0},
	}
	pm := NewManager(100.0, p)

	cfg := pm.GetCurrentConfig()
	if cfg.MovementSpeed != 80.0 {
		t.Errorf("expected speed 80, got %f", cfg.MovementSpeed)
	}

	pm.Update(50.0) // Transition to phase 1

	cfg = pm.GetCurrentConfig()
	if cfg.MovementSpeed != 100.0 {
		t.Errorf("expected speed 100, got %f", cfg.MovementSpeed)
	}
}
