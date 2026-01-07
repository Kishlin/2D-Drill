package input

import (
	"testing"
)

func TestInputState_HasMovementInput_LeftInput(t *testing.T) {
	inputState := InputState{Left: true}
	if !inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return true for Left input")
	}
}

func TestInputState_HasMovementInput_RightInput(t *testing.T) {
	inputState := InputState{Right: true}
	if !inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return true for Right input")
	}
}

func TestInputState_HasMovementInput_UpInput(t *testing.T) {
	inputState := InputState{Up: true}
	if !inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return true for Up input")
	}
}

func TestInputState_HasMovementInput_DrillInput(t *testing.T) {
	inputState := InputState{Drill: true}
	if !inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return true for Drill input")
	}
}

func TestInputState_HasMovementInput_SellInputOnly(t *testing.T) {
	inputState := InputState{Sell: true}
	if inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return false for Sell input only")
	}
}

func TestInputState_HasMovementInput_NoInput(t *testing.T) {
	inputState := InputState{}
	if inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return false for no inputs")
	}
}

func TestInputState_HasMovementInput_MultipleMovementInputs(t *testing.T) {
	inputState := InputState{Left: true, Up: true}
	if !inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return true for multiple movement inputs")
	}
}

func TestInputState_HasMovementInput_MovementWithSell(t *testing.T) {
	inputState := InputState{Right: true, Sell: true}
	if !inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return true when movement + sell input")
	}
}

func TestInputState_HasMovementInput_AllInputs(t *testing.T) {
	inputState := InputState{
		Left:  true,
		Right: true,
		Up:    true,
		Drill: true,
		Sell:  true,
	}
	if !inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return true with all inputs pressed")
	}
}

func TestInputState_HasMovementInput_DrillWithSell(t *testing.T) {
	inputState := InputState{Drill: true, Sell: true}
	if !inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return true for Drill + Sell")
	}
}

func TestInputState_HasMovementInput_SellWithNoMovementInputs(t *testing.T) {
	// Verify that Sell alone is NOT considered active
	inputState := InputState{Sell: true}
	expected := false
	actual := inputState.HasMovementInput()

	if actual != expected {
		t.Errorf("Sell-only input: expected HasMovementInput()=%v, got %v", expected, actual)
	}
}

func TestInputState_HasMovementInput_LeftOnly(t *testing.T) {
	inputState := InputState{Left: true}
	expected := true
	actual := inputState.HasMovementInput()

	if actual != expected {
		t.Errorf("Left-only input: expected HasMovementInput()=%v, got %v", expected, actual)
	}
}

func TestInputState_HasMovementInput_DrillAndLeftButNoOthers(t *testing.T) {
	inputState := InputState{Drill: true, Left: true}
	expected := true
	actual := inputState.HasMovementInput()

	if actual != expected {
		t.Errorf("Drill + Left: expected HasMovementInput()=%v, got %v", expected, actual)
	}
}
