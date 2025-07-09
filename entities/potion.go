package entities

import (
	"image"
	"rpg-go/camera"
	"rpg-go/constants"
	"rpg-go/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Potion struct {
	*Sprite
	AmtHeal uint
}

func (p *Potion) GetY() float64 {
	return p.Y
}

func (p *Potion) Draw(screen *ebiten.Image, cam *camera.Camera, _ *spritesheet.SpriteSheet) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	opts.GeoM.Translate(cam.X, cam.Y)

	frameRect := image.Rect(0, 0, constants.Tilesize, constants.Tilesize)

	screen.DrawImage(p.Img.SubImage(frameRect).(*ebiten.Image), opts)
}
