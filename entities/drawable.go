package entities

import (
	"rpg-go/camera"
	"rpg-go/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Drawable interface {
	GetY() float64
	Draw(screen *ebiten.Image, cam *camera.Camera, sheet *spritesheet.SpriteSheet)
}
