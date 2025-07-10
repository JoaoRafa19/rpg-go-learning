package spritesheet

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Assets struct {
	SkeletonImg *ebiten.Image
	PotionImg   *ebiten.Image
	DummyImg    *ebiten.Image
}

func LoadAssets() (*Assets, error) {
	skeletonImg, _, err := ebitenutil.NewImageFromFile("./assets/images/skeleton.png")
	if err != nil {
		return nil, err
	}
	dummyImg, _, err := ebitenutil.NewImageFromFile("./assets/images/dummy.png")
	if err != nil {
		return nil, err
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("./assets/images/health.png")
	if err != nil {
		return nil, err
	}

	return &Assets{skeletonImg, potionImg, dummyImg}, nil
}
