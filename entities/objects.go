package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"rpg-go/camera"
	"rpg-go/spritesheet"
)

type Objects struct {
	img        *ebiten.Image
	x, y, w, h float64
}

func NewObjects(img *ebiten.Image, x, y float64) *Objects {
	return &Objects{
		img: img,
		x:   x,
		y:   y,
	}
}

func (o *Objects) GetY() float64 {
	return o.y
}
func (o *Objects) Draw(screen *ebiten.Image, cam *camera.Camera, sheet *spritesheet.SpriteSheet) {
	opts := &ebiten.DrawImageOptions{}

	opts.GeoM.Reset()
	opts.GeoM.Translate(o.x, o.y)
	opts.GeoM.Translate(cam.X, cam.Y)

	screen.DrawImage(o.img, opts)
}
