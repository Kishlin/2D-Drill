package input

import (
	"testing"
)

func TestInputState_HasMovementInput_LeftInput(t *testing.T) {
	inputState := InputState{Left: true}
	if inputState.HasMovementInput() == false {
		t.Error("expected HasMovementInput() to return true for Left input")
	}
}

func TestInputState_HasMovementInput_RightInput(t *testing.T) {
	inputState := InputState{Right: true}
	if inputState.HasMovementInput() == false {
		t.Error("expected HasMovementInput() to return true for Right input")
	}
}

func TestInputState_HasMovementInput_UpInput(t *testing.T) {
	inputState := InputState{Up: true}
	if inputState.HasMovementInput() == false {
		t.Error("expected HasMovementInput() to return true for Up input")
	}
}

func TestInputState_HasMovementInput_DrillInput(t *testing.T) {
	inputState := InputState{Drill: true}
	if inputState.HasMovementInput() == false {
		t.Error("expected HasMovementInput() to return true for Drill input")
	}
}

func TestInputState_HasMovementInput_InteractInputOnly(t *testing.T) {
	inputState := InputState{Interact: true}
	if inputState.HasMovementInput() {
		t.Error("expected HasMovementInput() to return false for Interact input only")
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
	if inputState.HasMovementInput() == false {
		t.Error("expected HasMovementInput() to return true for multiple movement inputs")
	}
}

func TestInputState_HasMovementInput_MovementWithInteract(t *testing.T) {
	inputState := InputState{Right: true, Interact: true}
	if inputState.HasMovementInput() == false {
		t.Error("expected HasMovementInput() to return true when movement + sell input")
	}
}

func TestInputState_HasMovementInput_AllInputs(t *testing.T) {
	inputState := InputState{
		Left:  true,
		Right: true,
		Up:    true,
		Drill: true,
		Interact:  true,
	}
	if inputState.HasMovementInput() == false {
		t.Error("expected HasMovementInput() to return true with all inputs pressed")
	}
}

func TestInputState_HasMovementInput_DrillWithInteract(t *testing.T) {
	inputState := InputState{Drill: true, Interact: true}
	if inputState.HasMovementInput() == false {
		t.Error("expected HasMovementInput() to return true for Drill + Interact")
	}
}

func TestInputState_HasMovementInput_InteractWithNoMovementInputs(t *testing.T) {
	// Verify that Interact alone is NOT considered active
	inputState := InputState{Interact: true}
	expected := false
	actual := inputState.HasMovementInput()

	if actual != expected {
		t.Errorf("Interact-only input: expected HasMovementInput()=%v, got %v", expected, actual)
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
