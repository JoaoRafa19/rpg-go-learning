package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

func (s *Sprite) Width() float64 {
	return float64(s.Img.Bounds().Dx())
}
func (s *Sprite) Height() float64 {
	return float64(s.Img.Bounds().Dy())
}

func (s *Sprite) IsColliding(other *Sprite) bool {
	return other.X <= s.X+s.Width() && other.X > s.X && other.Y <= s.Y+s.Height() && other.Y > s.Y
}
