package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Img          *ebiten.Image
	X, Y, Dx, Dy float64
}

func (s *Sprite) Width() float64 {
	return float64(s.Img.Bounds().Dx())
}
func (s *Sprite) Height() float64 {
	return float64(s.Img.Bounds().Dy())
}
