package components

// Damageable is HP data for entities that can take damage.
// Vulnerability logic is entity-specific and not handled here.
type Damageable struct {
	HP    float32
	MaxHP float32
}

func NewDamageable(hp, maxHP float32) Damageable {
	return Damageable{
		HP:    hp,
		MaxHP: maxHP,
	}
}

func (d *Damageable) TakeDamage(amount float32) {
	d.HP -= amount
	if d.HP < 0 {
		d.HP = 0
	}
}

func (d *Damageable) IsAlive() bool    { return d.HP > 0 }
func (d *Damageable) IsDefeated() bool { return d.HP <= 0 }
