package entities

import "github.com/hajimehoshi/ebiten/v2"

type Projectile struct {
	*Sprite
	SpeedX, SpeedY float64
	Damage         int
	LifeSpan       float64 // qantos ticks o projetil dura
}

func NewProjectile(x, y, speedX, speedY float64, damage int, img *ebiten.Image) *Projectile {
	return &Projectile{
		Sprite:   &Sprite{X: x, Y: y, Img: img},
		SpeedX:   speedX,
		SpeedY:   speedY,
		Damage:   damage,
		LifeSpan: 120, // Vive por 2 segundos (120 ticks)
	}
}

func (p *Projectile) Update() {
	p.X += p.SpeedX
	p.Y += p.SpeedY
	p.LifeSpan--
}

func (p *Projectile) GetY() float64 {
	return p.Y
}

func (p *Projectile) Draw(screen *ebiten.Image) {

}
