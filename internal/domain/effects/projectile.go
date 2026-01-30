package effects

// ProjectileDamage applies damage from a projectile hit
type ProjectileDamage struct {
	Damage float32
}

func (e ProjectileDamage) Apply(ctx *EffectContext) {
	ctx.Player.DealDamage(e.Damage)
}
