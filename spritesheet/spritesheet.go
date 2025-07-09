package spritesheet

import "image"

type SpriteSheet struct {
	Width    int
	Height   int
	TileSize int
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := (index % s.Width) * s.TileSize
	y := (index / s.Width) * s.TileSize

	return image.Rect(x, y, x+s.TileSize, y+s.TileSize)
}

func NewSpriteSheet(Width, Height, TileSize int) *SpriteSheet {
	return &SpriteSheet{
		Width,
		Height,
		TileSize,
	}
}
