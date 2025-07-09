package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"rpg-go/camera"
	"rpg-go/spritesheet"
)

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

func (p *Projectile) Draw(screen *ebiten.Image, cam *camera.Camera, _ *spritesheet.SpriteSheet) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	opts.GeoM.Translate(cam.X, cam.Y)
	opts.GeoM.Rotate(p.LifeSpan)

	// Por enquanto, a shuriken não tem animação, então desenhamos a imagem inteira.
	// O ideal é que a imagem já seja do tamanho correto (ex: 16x16).
	screen.DrawImage(p.Img, opts)
}
