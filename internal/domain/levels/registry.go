package levels

import (
	"fmt"

	"github.com/Kishlin/drill-game/internal/domain/config"
)

// GetLevelConfig returns the game configuration for the specified level number
// Returns an error if the level doesn't exist
func GetLevelConfig(levelNum int) (*config.GameConfig, error) {
	switch levelNum {
	case -1:
		return GetTestLevelConfig(), nil
	case 1:
		return GetLevel1Config(), nil
	default:
		return nil, fmt.Errorf("level %d not found", levelNum)
	}
}
