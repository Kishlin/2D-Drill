package effects

type SetFuel struct {
	Amount float32
}

func (e SetFuel) Apply(ctx *EffectContext) {
	ctx.Player.Fuel = e.Amount
}

type SetHP struct {
	Amount float32
}

func (e SetHP) Apply(ctx *EffectContext) {
	ctx.Player.HP = e.Amount
}
