package main

import (
	"image"
	"log"
	"rpg-go/constants"
	"rpg-go/entities"

	"github.com/hajimehoshi/ebiten/v2"
)

func CheckCollisionsVerticaly(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+constants.Tilesize, int(sprite.Y)+constants.Tilesize)) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - constants.Tilesize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}
func CheckCollisionsHorizontaly(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+constants.Tilesize, int(sprite.Y)+constants.Tilesize)) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - constants.Tilesize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("RPG Go!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
