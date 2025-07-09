package entities

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"rpg-go/camera"
	"rpg-go/constants"
	"rpg-go/spritesheet"
)

type TrainingDummy struct {
	*Sprite
	IsAnimating   bool
	AnimationTick int
	currentFrame  int
	framecount    int
}

func NewTrainingDummy(x, y float64, img *ebiten.Image) *TrainingDummy {
	return &TrainingDummy{
		IsAnimating:   false,
		AnimationTick: 0,
		currentFrame:  1,
		framecount:    4,
		Sprite: &Sprite{
			Img: img,
			X:   x,
			Y:   y,
		},
	}
}

func (t *TrainingDummy) Update() {
	if !t.IsAnimating {
		return
	}

	t.AnimationTick++

	if t.AnimationTick > 10 {
		t.AnimationTick = 0

		t.currentFrame = (t.currentFrame + 1) % t.framecount

	}
	if t.AnimationTick == 0 && t.currentFrame == 0 {
		t.IsAnimating = false
	}
}

func (d *TrainingDummy) Hit() {
	if !d.IsAnimating {
		d.IsAnimating = true
		d.AnimationTick = 0
		d.currentFrame = 1 // Começa a animação no segundo frame
	}

	fmt.Println("hit")
}

// Drawable
func (d *TrainingDummy) GetY() float64 {
	return d.Y
}

func (d *TrainingDummy) Draw(screen *ebiten.Image, cam *camera.Camera, _ *spritesheet.SpriteSheet) {
	frameRect := image.Rect(
		d.currentFrame*16, // 0*32 ou 1*32
		0,
		(d.currentFrame+1)*constants.Tilesize,
		constants.Tilesize*2,
	)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(d.X, d.Y)
	opts.GeoM.Translate(cam.X, cam.Y)

	screen.DrawImage(d.Img.SubImage(frameRect).(*ebiten.Image), opts)
}
