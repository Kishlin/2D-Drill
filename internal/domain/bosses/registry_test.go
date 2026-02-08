package bosses

import (
	"strings"
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/bosses/phases"
	"github.com/Kishlin/drill-game/internal/domain/bosses/statemachine"
	"github.com/Kishlin/drill-game/internal/domain/entities"
	"github.com/Kishlin/drill-game/internal/domain/projectiles"
	"github.com/Kishlin/drill-game/internal/domain/types"
)

// saveRegistry snapshots the registry map and restores it via t.Cleanup.
func saveRegistry(t *testing.T) {
	t.Helper()
	original := make(map[string]BossConstructor, len(registry))
	for k, v := range registry {
		original[k] = v
	}
	t.Cleanup(func() {
		registry = original
	})
}

// stubBoss is a minimal Boss implementation for registry tests.
type stubBoss struct {
	*BaseBoss
}

func (s *stubBoss) Update(_ *entities.Player, _ float32) []projectiles.SpawnRequest {
	return nil
}

func newStubBoss(roomStartY, worldWidth float32) *stubBoss {
	base := NewBaseBoss(BaseBossConfig{
		Position: types.NewVec2(worldWidth/2, roomStartY),
		MaxHP:    10,
		BoxSet:   NewBoxSet(nil, nil, nil),
		Phases:   []phases.Config{{HPThreshold: 0}},
	})
	b := &stubBoss{BaseBoss: base}
	b.Self = b
	b.SetStateMachine(statemachine.NewStateMachine(
		map[statemachine.StateID]*statemachine.State{
			0: {ID: 0, OnUpdate: func(ctx *statemachine.StateContext) statemachine.StateResult {
				return statemachine.StateResult{NextState: statemachine.StateIDNone}
			}},
		}, 0,
	))
	return b
}

func TestRegister_AndCreate_HappyPath(t *testing.T) {
	saveRegistry(t)

	Register("stub_boss", func(roomStartY, worldWidth float32) Boss {
		return newStubBoss(roomStartY, worldWidth)
	})

	boss, err := Create("stub_boss", 100, 500)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if boss == nil {
		t.Fatal("Expected non-nil boss")
	}
}

func TestCreate_UnknownType_ReturnsError(t *testing.T) {
	saveRegistry(t)

	boss, err := Create("nonexistent_boss", 100, 500)
	if err == nil {
		t.Fatal("Expected error for unknown boss type")
	}
	if strings.Contains(err.Error(), "unknown boss type") == false {
		t.Errorf("Error should contain 'unknown boss type', got %q", err.Error())
	}
	if boss != nil {
		t.Error("Expected nil boss on error")
	}
}

func TestRegister_DuplicatePanics(t *testing.T) {
	saveRegistry(t)

	constructor := func(roomStartY, worldWidth float32) Boss {
		return newStubBoss(roomStartY, worldWidth)
	}

	Register("dup_boss", constructor)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic on duplicate registration")
		}
		msg, ok := r.(string)
		if ok == false {
			t.Fatalf("Expected string panic, got %T", r)
		}
		if strings.Contains(msg, "dup_boss") == false {
			t.Errorf("Panic message should contain boss type, got %q", msg)
		}
	}()

	Register("dup_boss", constructor)
}

func TestCreate_PassesParametersToConstructor(t *testing.T) {
	saveRegistry(t)

	var capturedRoomStartY, capturedWorldWidth float32

	Register("param_boss", func(roomStartY, worldWidth float32) Boss {
		capturedRoomStartY = roomStartY
		capturedWorldWidth = worldWidth
		return newStubBoss(roomStartY, worldWidth)
	})

	_, err := Create("param_boss", 42, 1024)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if capturedRoomStartY != 42 {
		t.Errorf("Expected roomStartY 42, got %f", capturedRoomStartY)
	}
	if capturedWorldWidth != 1024 {
		t.Errorf("Expected worldWidth 1024, got %f", capturedWorldWidth)
	}
}
