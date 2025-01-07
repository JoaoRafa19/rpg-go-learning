package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type StartScene struct {
	loaded bool
}

func NewStartScene() *StartScene {
	return &StartScene{
		loaded: false,
	}
}

func (s *StartScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	image := ebiten.NewImage(200, 100)
	opts := ebiten.DrawImageOptions{}

	opts.GeoM.Translate(100, 100)

	ebitenutil.DebugPrint(image, "Press ENTER to start.")
	screen.DrawImage(image, &opts)
	opts.GeoM.Reset()

}

func (s *StartScene) FirstLoad() {
	s.loaded = true
}

func (s *StartScene) IsLoaded() bool {
	return s.loaded
}

func (s *StartScene) OnEnter() {
}

func (s *StartScene) OnExit() {
}

func (s *StartScene) Update() SceneId {

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	return StartSceneId
}

var _ Scene = (*StartScene)(nil)
