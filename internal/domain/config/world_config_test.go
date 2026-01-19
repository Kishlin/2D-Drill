package config_test

import (
	"testing"

	"github.com/Kishlin/drill-game/internal/domain/config"
)

func TestWorldConfig_Validate_ValidConfig(t *testing.T) {
	cfg := validTestCfg()
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Valid config should not return error, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativeWidth(t *testing.T) {
	cfg := validTestCfg()
	cfg.Width = -100
	err := cfg.Validate()
	if err == nil {
		t.Error("Negative width should return error")
	}
}

func TestWorldConfig_Validate_ZeroWidth(t *testing.T) {
	cfg := validTestCfg()
	cfg.Width = 0
	err := cfg.Validate()
	if err == nil {
		t.Error("Zero width should return error")
	}
}

func TestWorldConfig_Validate_NegativeHeight(t *testing.T) {
	cfg := validTestCfg()
	cfg.Height = -100
	err := cfg.Validate()
	if err == nil {
		t.Error("Negative height should return error")
	}
}

func TestWorldConfig_Validate_ZeroHeight(t *testing.T) {
	cfg := validTestCfg()
	cfg.Height = 0
	err := cfg.Validate()
	if err == nil {
		t.Error("Zero height should return error")
	}
}

func TestWorldConfig_Validate_NegativeGroundLevel(t *testing.T) {
	cfg := validTestCfg()
	cfg.GroundLevel = -10
	err := cfg.Validate()
	if err == nil {
		t.Error("Negative ground level should return error")
	}
}

func TestWorldConfig_Validate_GroundLevelAboveHeight(t *testing.T) {
	cfg := validTestCfg()
	cfg.GroundLevel = cfg.Height + 100
	err := cfg.Validate()
	if err == nil {
		t.Error("Ground level > height should return error")
	}
}

func TestWorldConfig_Validate_GroundLevelAtZero(t *testing.T) {
	cfg := validTestCfg()
	cfg.GroundLevel = 0
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Ground level at 0 should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_GroundLevelAtHeight(t *testing.T) {
	cfg := validTestCfg()
	cfg.GroundLevel = cfg.Height
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Ground level at height should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativePlayerSpawnX(t *testing.T) {
	cfg := validTestCfg()
	cfg.PlayerSpawn.X = -100
	err := cfg.Validate()
	if err == nil {
		t.Error("Negative player spawn X should return error")
	}
}

func TestWorldConfig_Validate_PlayerSpawnXBeyondWidth(t *testing.T) {
	cfg := validTestCfg()
	cfg.PlayerSpawn.X = cfg.Width + 100
	err := cfg.Validate()
	if err == nil {
		t.Error("Player spawn X > width should return error")
	}
}

func TestWorldConfig_Validate_PlayerSpawnXAtZero(t *testing.T) {
	cfg := validTestCfg()
	cfg.PlayerSpawn.X = 0
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Player spawn X at 0 should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_PlayerSpawnXAtWidth(t *testing.T) {
	cfg := validTestCfg()
	cfg.PlayerSpawn.X = cfg.Width
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Player spawn X at width should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativePlayerSpawnY(t *testing.T) {
	cfg := validTestCfg()
	cfg.PlayerSpawn.Y = -100
	err := cfg.Validate()
	if err == nil {
		t.Error("Negative player spawn Y should return error")
	}
}

func TestWorldConfig_Validate_PlayerSpawnYBeyondHeight(t *testing.T) {
	cfg := validTestCfg()
	cfg.PlayerSpawn.Y = cfg.Height + 100
	err := cfg.Validate()
	if err == nil {
		t.Error("Player spawn Y > height should return error")
	}
}

func TestWorldConfig_Validate_PlayerSpawnYAtZero(t *testing.T) {
	cfg := validTestCfg()
	cfg.PlayerSpawn.Y = 0
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Player spawn Y at 0 should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_PlayerSpawnYAtHeight(t *testing.T) {
	cfg := validTestCfg()
	cfg.PlayerSpawn.Y = cfg.Height
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Player spawn Y at height should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativeMarketX_Partial(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.MarketX = -100 // Partially off-screen (valid)
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Market partially off-screen left should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativeMarketX_CompletelyOffScreen(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.MarketX = -500 // Completely off-screen (invalid, -500 + 320 = -180)
	err := cfg.Validate()
	if err == nil {
		t.Error("Market completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_MarketXExtendsRightBeyondWidth_Partial(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.MarketX = cfg.Width - 100 // Partially off-screen right (valid)
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Market partially off-screen right should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_MarketXAtZero(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.MarketX = 0
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Market X at 0 should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_MarketXAtRightEdge(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.MarketX = cfg.Width - 320 // Exactly fits
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Market X at right edge should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_NegativeFuelStationX_CompletelyOffScreen(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.FuelStationX = -500 // Completely off-screen
	err := cfg.Validate()
	if err == nil {
		t.Error("Fuel station completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_NegativeHospitalX_CompletelyOffScreen(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.HospitalX = -500 // Completely off-screen
	err := cfg.Validate()
	if err == nil {
		t.Error("Hospital completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_UpgradeShopXExtendsRight_Partial(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.UpgradeShopX = cfg.Width - 100 // Partially off-screen right (valid)
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Upgrade shop partially off-screen right should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_UpgradeShopXExtendsRight_Completely(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.UpgradeShopX = cfg.Width + 100 // Completely off-screen
	err := cfg.Validate()
	if err == nil {
		t.Error("Upgrade shop completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_ItemShopXExtendsRight_Partial(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.ItemShopX = cfg.Width - 100 // Partially off-screen right (valid)
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Item shop partially off-screen right should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_ItemShopXExtendsRight_Completely(t *testing.T) {
	cfg := validTestCfg()
	cfg.BuildingLayout.ItemShopX = cfg.Width + 100 // Completely off-screen
	err := cfg.Validate()
	if err == nil {
		t.Error("Item shop completely off-screen should return error")
	}
}

func TestWorldConfig_Validate_AllBuildingsAtSameLocation(t *testing.T) {
	cfg := validTestCfg()
	// Place all buildings at center - valid positions
	centerX := cfg.Width/2 - 160 // Center the 320-wide buildings
	cfg.BuildingLayout.MarketX = centerX
	cfg.BuildingLayout.FuelStationX = centerX
	cfg.BuildingLayout.HospitalX = centerX
	cfg.BuildingLayout.UpgradeShopX = centerX
	cfg.BuildingLayout.ItemShopX = centerX
	err := cfg.Validate()
	if err != nil {
		t.Errorf("All buildings at same location should be valid, got: %v", err)
	}
}

func TestWorldConfig_Validate_SmallWorld(t *testing.T) {
	cfg := &config.WorldConfig{
		Width:       640,
		Height:      1280,
		GroundLevel: 320,
		Seed:        1,
		PlayerSpawn: config.PlayerSpawn{X: 160, Y: 310},
		BuildingLayout: config.BuildingLayout{
			MarketX:      100,
			FuelStationX: 0,
			HospitalX:    200,
			UpgradeShopX: 300,
			ItemShopX:    0,
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Small valid world should validate, got: %v", err)
	}
}

func TestWorldConfig_Validate_LargeWorld(t *testing.T) {
	cfg := &config.WorldConfig{
		Width:       1000000,
		Height:      10000000,
		GroundLevel: 500000,
		Seed:        999,
		PlayerSpawn: config.PlayerSpawn{X: 500000, Y: 499990},
		BuildingLayout: config.BuildingLayout{
			MarketX:      100000,
			FuelStationX: 50000,
			HospitalX:    200000,
			UpgradeShopX: 300000,
			ItemShopX:    400000,
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Large valid world should validate, got: %v", err)
	}
}

func TestWorldConfig_Validate_CompactLayout(t *testing.T) {
	// This is the actual configuration from main.go with compact building layout
	width := float32(2360.0)
	cfg := &config.WorldConfig{
		Width:       width,
		Height:      64000,
		GroundLevel: 640,
		Seed:        42,
		PlayerSpawn: config.PlayerSpawn{X: width / 2, Y: 570},
		BuildingLayout: config.BuildingLayout{
			HospitalX:    100,
			FuelStationX: 460,
			MarketX:      820,
			UpgradeShopX: 1180,
			ItemShopX:    1540,
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("Compact layout should be valid, got: %v", err)
	}
}

// Test helpers

func validTestCfg() *config.WorldConfig {
	return &config.WorldConfig{
		Width:       7680,
		Height:      64000,
		GroundLevel: 640,
		Seed:        42,
		PlayerSpawn: config.PlayerSpawn{X: 3840, Y: 630},
		BuildingLayout: config.BuildingLayout{
			MarketX:      840,
			FuelStationX: 120,
			HospitalX:    0, // Start at left edge (valid)
			UpgradeShopX: 1290,
			ItemShopX:    1650,
		},
	}
}
