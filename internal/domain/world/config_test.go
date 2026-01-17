package world

import (
	"testing"
)

// Valid config factory for tests
func validTestConfig() *WorldConfig {
	return &WorldConfig{
		Width:       7680,
		Height:      64000,
		GroundLevel: 640,
		Seed:        42,
		PlayerSpawn: PlayerSpawn{X: 3840, Y: 630},
		BuildingLayout: BuildingLayout{
			MarketX:      840,
			FuelStationX: 120,
			HospitalX:    0,      // Start at left edge (valid)
			UpgradeShopX: 1290,
			ItemShopX:    1650,
		},
	}
}

func TestWorldConfig_Validate_ValidConfig(t *testing.T) {
	config := validTestConfig()
	err := config.Validate()
	if err != nil {
		t.Errorf("Valid config should not return error, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativeWidth(t *testing.T) {
	config := validTestConfig()
	config.Width = -100
	err := config.Validate()
	if err == nil {
		t.Error("Negative width should return error")
	}
}

func TestWorldConfig_Validate_ZeroWidth(t *testing.T) {
	config := validTestConfig()
	config.Width = 0
	err := config.Validate()
	if err == nil {
		t.Error("Zero width should return error")
	}
}

func TestWorldConfig_Validate_NegativeHeight(t *testing.T) {
	config := validTestConfig()
	config.Height = -100
	err := config.Validate()
	if err == nil {
		t.Error("Negative height should return error")
	}
}

func TestWorldConfig_Validate_ZeroHeight(t *testing.T) {
	config := validTestConfig()
	config.Height = 0
	err := config.Validate()
	if err == nil {
		t.Error("Zero height should return error")
	}
}

func TestWorldConfig_Validate_NegativeGroundLevel(t *testing.T) {
	config := validTestConfig()
	config.GroundLevel = -10
	err := config.Validate()
	if err == nil {
		t.Error("Negative ground level should return error")
	}
}

func TestWorldConfig_Validate_GroundLevelAboveHeight(t *testing.T) {
	config := validTestConfig()
	config.GroundLevel = config.Height + 100
	err := config.Validate()
	if err == nil {
		t.Error("Ground level > height should return error")
	}
}

func TestWorldConfig_Validate_GroundLevelAtZero(t *testing.T) {
	config := validTestConfig()
	config.GroundLevel = 0
	err := config.Validate()
	if err != nil {
		t.Errorf("Ground level at 0 should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_GroundLevelAtHeight(t *testing.T) {
	config := validTestConfig()
	config.GroundLevel = config.Height
	err := config.Validate()
	if err != nil {
		t.Errorf("Ground level at height should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativePlayerSpawnX(t *testing.T) {
	config := validTestConfig()
	config.PlayerSpawn.X = -100
	err := config.Validate()
	if err == nil {
		t.Error("Negative player spawn X should return error")
	}
}

func TestWorldConfig_Validate_PlayerSpawnXBeyondWidth(t *testing.T) {
	config := validTestConfig()
	config.PlayerSpawn.X = config.Width + 100
	err := config.Validate()
	if err == nil {
		t.Error("Player spawn X > width should return error")
	}
}

func TestWorldConfig_Validate_PlayerSpawnXAtZero(t *testing.T) {
	config := validTestConfig()
	config.PlayerSpawn.X = 0
	err := config.Validate()
	if err != nil {
		t.Errorf("Player spawn X at 0 should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_PlayerSpawnXAtWidth(t *testing.T) {
	config := validTestConfig()
	config.PlayerSpawn.X = config.Width
	err := config.Validate()
	if err != nil {
		t.Errorf("Player spawn X at width should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativePlayerSpawnY(t *testing.T) {
	config := validTestConfig()
	config.PlayerSpawn.Y = -100
	err := config.Validate()
	if err == nil {
		t.Error("Negative player spawn Y should return error")
	}
}

func TestWorldConfig_Validate_PlayerSpawnYBeyondHeight(t *testing.T) {
	config := validTestConfig()
	config.PlayerSpawn.Y = config.Height + 100
	err := config.Validate()
	if err == nil {
		t.Error("Player spawn Y > height should return error")
	}
}

func TestWorldConfig_Validate_PlayerSpawnYAtZero(t *testing.T) {
	config := validTestConfig()
	config.PlayerSpawn.Y = 0
	err := config.Validate()
	if err != nil {
		t.Errorf("Player spawn Y at 0 should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_PlayerSpawnYAtHeight(t *testing.T) {
	config := validTestConfig()
	config.PlayerSpawn.Y = config.Height
	err := config.Validate()
	if err != nil {
		t.Errorf("Player spawn Y at height should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativeMarketX_Partial(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.MarketX = -100 // Partially off-screen (valid)
	err := config.Validate()
	if err != nil {
		t.Errorf("Market partially off-screen left should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativeMarketX_CompletelyOffScreen(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.MarketX = -500 // Completely off-screen (invalid, -500 + 320 = -180)
	err := config.Validate()
	if err == nil {
		t.Error("Market completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_MarketXExtendsRightBeyondWidth_Partial(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.MarketX = config.Width - 100 // Partially off-screen right (valid)
	err := config.Validate()
	if err != nil {
		t.Errorf("Market partially off-screen right should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_MarketXAtZero(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.MarketX = 0
	err := config.Validate()
	if err != nil {
		t.Errorf("Market X at 0 should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_MarketXAtRightEdge(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.MarketX = config.Width - 320 // Exactly fits
	err := config.Validate()
	if err != nil {
		t.Errorf("Market X at right edge should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativeFuelStationX_CompletelyOffScreen(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.FuelStationX = -500 // Completely off-screen
	err := config.Validate()
	if err == nil {
		t.Error("Fuel station completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_NegativeHospitalX_CompletelyOffScreen(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.HospitalX = -500 // Completely off-screen
	err := config.Validate()
	if err == nil {
		t.Error("Hospital completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_UpgradeShopXExtendsRight_Partial(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.UpgradeShopX = config.Width - 100 // Partially off-screen right (valid)
	err := config.Validate()
	if err != nil {
		t.Errorf("Upgrade shop partially off-screen right should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_UpgradeShopXExtendsRight_Completely(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.UpgradeShopX = config.Width + 100 // Completely off-screen
	err := config.Validate()
	if err == nil {
		t.Error("Upgrade shop completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_ItemShopXExtendsRight_Partial(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.ItemShopX = config.Width - 100 // Partially off-screen right (valid)
	err := config.Validate()
	if err != nil {
		t.Errorf("Item shop partially off-screen right should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_ItemShopXExtendsRight_Completely(t *testing.T) {
	config := validTestConfig()
	config.BuildingLayout.ItemShopX = config.Width + 100 // Completely off-screen
	err := config.Validate()
	if err == nil {
		t.Error("Item shop completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_AllBuildingsAtSameLocation(t *testing.T) {
	config := validTestConfig()
	// Place all buildings at center - valid positions
	centerX := config.Width/2 - 160 // Center the 320-wide buildings
	config.BuildingLayout.MarketX = centerX
	config.BuildingLayout.FuelStationX = centerX
	config.BuildingLayout.HospitalX = centerX
	config.BuildingLayout.UpgradeShopX = centerX
	config.BuildingLayout.ItemShopX = centerX
	err := config.Validate()
	if err != nil {
		t.Errorf("All buildings at same location should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_SmallWorld(t *testing.T) {
	config := &WorldConfig{
		Width:       640,
		Height:      1280,
		GroundLevel: 320,
		Seed:        1,
		PlayerSpawn: PlayerSpawn{X: 160, Y: 310},
		BuildingLayout: BuildingLayout{
			MarketX:      100,
			FuelStationX: 0,
			HospitalX:    200,
			UpgradeShopX: 300,
			ItemShopX:    0,
		},
	}
	err := config.Validate()
	if err != nil {
		t.Errorf("Small valid world should validate, got: %v", err)
	}
}

func TestWorldConfig_Validate_LargeWorld(t *testing.T) {
	config := &WorldConfig{
		Width:       1000000,
		Height:      10000000,
		GroundLevel: 500000,
		Seed:        999,
		PlayerSpawn: PlayerSpawn{X: 500000, Y: 499990},
		BuildingLayout: BuildingLayout{
			MarketX:      100000,
			FuelStationX: 50000,
			HospitalX:    200000,
			UpgradeShopX: 300000,
			ItemShopX:    400000,
		},
	}
	err := config.Validate()
	if err != nil {
		t.Errorf("Large valid world should validate, got: %v", err)
	}
}

func TestNewWorldConfigForTesting_IsValid(t *testing.T) {
	config := NewWorldConfigForTesting(1280, 720, 640, 42)
	err := config.Validate()
	if err != nil {
		t.Errorf("Testing config should be valid, got: %v", err)
	}
}

func TestNewWorldConfigForTesting_HasCorrectDimensions(t *testing.T) {
	config := NewWorldConfigForTesting(1280, 720, 640, 42)
	if config.Width != 1280 {
		t.Errorf("Expected width 1280, got %f", config.Width)
	}
	if config.Height != 720 {
		t.Errorf("Expected height 720, got %f", config.Height)
	}
	if config.GroundLevel != 640 {
		t.Errorf("Expected ground level 640, got %f", config.GroundLevel)
	}
	if config.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", config.Seed)
	}
}

func TestNewWorldConfigForTesting_CentersBuildingsAndPlayer(t *testing.T) {
	width := float32(1280)
	config := NewWorldConfigForTesting(width, 720, 640, 42)

	// All buildings should be centered
	expectedX := width / 2
	if config.BuildingLayout.MarketX != expectedX {
		t.Errorf("Expected market X %f, got %f", expectedX, config.BuildingLayout.MarketX)
	}
	if config.BuildingLayout.FuelStationX != expectedX {
		t.Errorf("Expected fuel station X %f, got %f", expectedX, config.BuildingLayout.FuelStationX)
	}

	// Player spawn should be centered horizontally, offset vertically
	if config.PlayerSpawn.X != expectedX {
		t.Errorf("Expected player spawn X %f, got %f", expectedX, config.PlayerSpawn.X)
	}
	if config.PlayerSpawn.Y != 630 {
		t.Errorf("Expected player spawn Y 630, got %f", config.PlayerSpawn.Y)
	}
}

func TestWorldConfig_Validate_CompactLayout(t *testing.T) {
	// This is the actual configuration from main.go with compact building layout
	width := float32(2360.0)
	config := &WorldConfig{
		Width:       width,
		Height:      64000,
		GroundLevel: 640,
		Seed:        42,
		PlayerSpawn: PlayerSpawn{X: width / 2, Y: 570},
		BuildingLayout: BuildingLayout{
			HospitalX:    100,
			FuelStationX: 460,
			MarketX:      820,
			UpgradeShopX: 1180,
			ItemShopX:    1540,
		},
	}
	err := config.Validate()
	if err != nil {
		t.Errorf("Compact layout should be valid, got: %v", err)
	}
}
