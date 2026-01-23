package effects

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Apply(ctx *EffectContext, effects []Effect) {
	for _, effect := range effects {
		effect.Apply(ctx)
	}
}
