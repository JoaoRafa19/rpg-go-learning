package components

type Combat interface {
	Health() int
	AttackPower() int
	Attacking() bool
	Attack() bool
	Update()
	Damage(amount int)
	MaxHealth() int
}

type BasicCombat struct {
	health      int
	attackPower int
	attacking   bool
	maxHeath    int
}

func (b *BasicCombat) Heal(i int) {
	b.health += i
	// if b.health >= b.maxHeath {
	// 	b.health = b.maxHeath
	// }
}

func NewBasicCombat(health, attackPower int) *BasicCombat {
	return &BasicCombat{
		health:      health,
		attackPower: attackPower,
		attacking:   false,
		maxHeath:    health,
	}
}

func (b *BasicCombat) MaxHealth() int {
	return b.maxHeath
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

func (b *BasicCombat) Attacking() bool {
	return b.attacking
}

func (b *BasicCombat) Attack() bool {
	b.attacking = true
	return true
}

func (b *BasicCombat) Update() {

}

var _ Combat = (*BasicCombat)(nil)

type EnemyCombat struct {
	*BasicCombat
	attackCooldown  int
	timeSinceAttack int
}

func NewEnemieCombat(health, attackPower, attackCooldown int) *EnemyCombat {
	return &EnemyCombat{
		BasicCombat:     NewBasicCombat(health, attackPower),
		attackCooldown:  attackCooldown,
		timeSinceAttack: 0,
	}
}

func (e *EnemyCombat) Attack() bool {
	if e.timeSinceAttack >= e.attackCooldown {
		e.attacking = true
		e.timeSinceAttack = 0
		return true
	}
	return false
}

func (e *EnemyCombat) Update() {
	e.timeSinceAttack += 1
}

var _ Combat = (*EnemyCombat)(nil)

type DummyCombat struct {
	*BasicCombat
}
