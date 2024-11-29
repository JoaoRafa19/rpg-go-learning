package components

type Combat interface {
	Health() int
	AttackPower() int
	Attackin() bool
	Attack()
	Update()
	Damage(amount int)
}

type BasicCombat struct {
	health      int
	attackPower int
	attacking   bool
}

func NewBasicCombat(health, attackPower int) *BasicCombat {
	return &BasicCombat{
		health,
		attackPower,
		false,
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

func (b *BasicCombat) Attackin() bool {
	return b.attacking
}

func (b *BasicCombat) Attack() {
	b.attacking = true
}

func (b *BasicCombat) Update() {

}

var _ Combat = (*BasicCombat)(nil)

type EnemyCombat struct {
	*BasicCombat
	attackCooldown int 
	timeSinceAttack int 

}

func NewEnemieCombat() {
	
}

func (e* EnemyCombat) Attack() {
	if e.timeSinceAttack >= e.attackCooldown {
		e.attacking = true
		e.timeSinceAttack = 0
	} 
}

func (e*EnemyCombat) Update() {
	e.timeSinceAttack += 1
}

var _ Combat = (*EnemyCombat)(nil)
