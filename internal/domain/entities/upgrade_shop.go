package entities

import "github.com/Kishlin/drill-game/internal/domain/types"

type UpgradeType int

const (
	UpgradeTypeEngine   UpgradeType = 0
	UpgradeTypeHull     UpgradeType = 1
	UpgradeTypeFuelTank UpgradeType = 2
)

const (
	UpgradeShopWidth  = 320.0
	UpgradeShopHeight = 192.0
)

type UpgradeShop struct {
	AABB        types.AABB
	UpgradeType UpgradeType
}

func NewUpgradeShop(x, y float32, upgradeType UpgradeType) *UpgradeShop {
	return &UpgradeShop{
		AABB:        types.NewAABB(x, y, UpgradeShopWidth, UpgradeShopHeight),
		UpgradeType: upgradeType,
	}
}

func (us *UpgradeShop) IsPlayerInRange(player *Player) bool {
	return us.AABB.Intersects(player.AABB)
}
