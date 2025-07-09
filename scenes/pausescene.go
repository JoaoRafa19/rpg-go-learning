package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	loaded bool
}

func NewPauseScene() *PauseScene {
	return &PauseScene{
		loaded: false,
	}
}

func (s *PauseScene) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{100, 100, 120, 100})

	titileImage := ebiten.NewImage(200, 100)
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(150, 10)
	ebitenutil.DebugPrint(titileImage, "PAUSED")
	screen.DrawImage(titileImage, &opts)
	opts.GeoM.Reset()

	image := ebiten.NewImage(200, 100)
	opts.GeoM.Translate(100, 100)
	ebitenutil.DebugPrint(image, "Press ENTER to unpause.")
	screen.DrawImage(image, &opts)
	opts.GeoM.Reset()

	exitImage := ebiten.NewImage(200, 100)
	opts.GeoM.Translate(100, 150)
	ebitenutil.DebugPrint(exitImage, "Press ESCAPE to exit.")
	screen.DrawImage(exitImage, &opts)
	opts.GeoM.Reset()
}

func (s *PauseScene) FirstLoad() {
	s.loaded = true
}

func (s *PauseScene) IsLoaded() bool {
	return s.loaded
}

func (s *PauseScene) OnEnter() {
}

func (s *PauseScene) OnExit() {
}

func (s *PauseScene) Update() SceneId {

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ExitSceneId
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	return PauseSceneId
}

var _ Scene = (*PauseScene)(nil)
