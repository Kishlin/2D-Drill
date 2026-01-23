package effects

type TakeMoney struct {
	Amount int
}

func (e TakeMoney) Apply(ctx *EffectContext) {
	ctx.Player.Money -= e.Amount
}

type AddMoney struct {
	Amount int
}

func (e AddMoney) Apply(ctx *EffectContext) {
	ctx.Player.Money += e.Amount
}
