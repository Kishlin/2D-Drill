package effects

// HazardDamage applies flat damage when drilling a hazard tile
type HazardDamage struct {
	Damage float32
}

func (e HazardDamage) Apply(ctx *EffectContext) {
	ctx.Player.DealDamage(e.Damage)
}

// HazardHeatDamage applies damage scaled by heat resistance when drilling a hazard tile
// Formula: finalDamage = baseDamage - (currentResistance / maxResistance) * baseDamage * maxDamageReduction
type HazardHeatDamage struct {
	BaseDamage         float32
	MaxHeatResistance  float32
	MaxDamageReduction float32
}

func (e HazardHeatDamage) Apply(ctx *EffectContext) {
	currentHeatResistance := ctx.Player.HeatResistance()

	// Linear reduction: 100% damage at 0 resistance, reduced by maxDamageReduction at max resistance
	damageReduction := (currentHeatResistance / e.MaxHeatResistance) * e.BaseDamage * e.MaxDamageReduction
	damage := e.BaseDamage - damageReduction

	if damage < 0 {
		damage = 0
	}

	ctx.Player.DealDamage(damage)
}

// HazardMoney gives money when drilling a hazard tile
type HazardMoney struct {
	Amount int
}

func (e HazardMoney) Apply(ctx *EffectContext) {
	ctx.Player.Money += e.Amount
}
