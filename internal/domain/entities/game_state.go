package entities

type GameState int

const (
	GameStatePlaying GameState = iota
	GameStateVictory
	GameStateDefeat
)
