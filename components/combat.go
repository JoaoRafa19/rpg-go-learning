package components

type Combat interface {
	Health() int
	AttackPower() int
	Damage(amount int)
}

type BasicCombat struct {
	health      int
	attackPower int
}

func NewBasicCombat(health, attackPower int) *BasicCombat {
	return &BasicCombat{
		health,
		attackPower,
	}
}

func (b *BasicCombat) Damage(amount int) {
	b.health -= amount
}

// AttackPower implements Combat.
func (b *BasicCombat) AttackPower() int {
	return b.attackPower
}

func (b *BasicCombat) Health() int {
	return b.health
}

var _ Combat = (*BasicCombat)(nil)
