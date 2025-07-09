package entities

import (
	"image"
	"rpg-go/camera"
	"rpg-go/components"
	"rpg-go/constants"
	"rpg-go/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	*Sprite
	FollowsPlayer bool
	CombatComp    *components.EnemyCombat
}

func (e *Enemy) GetY() float64 {
	return e.Y
}

func (e *Enemy) Draw(screen *ebiten.Image, cam *camera.Camera, sheet *spritesheet.SpriteSheet) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(e.X, e.Y)
	opts.GeoM.Translate(cam.X, cam.Y)

	// Define o retângulo para o primeiro sprite (assumindo que o inimigo ocupa 1 tile)
	frameRect := image.Rect(0, 0, constants.Tilesize, constants.Tilesize)

	// Desenha a sub-imagem
	screen.DrawImage(e.Img.SubImage(frameRect).(*ebiten.Image), opts)
}
