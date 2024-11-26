package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (s *Sprite) isColliding(other *Sprite) bool {
	return other.X <= s.X+s.Width() && other.X > s.X && other.Y <= s.Y+s.Height() && other.Y > s.Y
}

type Player struct {
	*Sprite
	Health uint
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Potion struct {
	*Sprite
	AmtHeal uint
}

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

type Game struct {
	Player       *Player
	sprites      []*Enemy
	potions      []*Potion
	TilemapJSON  *TilemapJSON
	TilemapImage *ebiten.Image
}

func (g *Game) Update() error {
	// React to keyPresses
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Player.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Player.X += 2
	}

	for _, sprite := range g.sprites {
		if sprite.FollowsPlayer {
			if sprite.X < g.Player.X {
				sprite.X += 1
			} else if sprite.X > g.Player.X {
				sprite.X -= 1
			}
			if sprite.Y < g.Player.Y {
				sprite.Y += 1
			} else if sprite.Y > g.Player.Y {
				sprite.Y -= 1
			}
		}
	}

	for i, potion := range g.potions {
		// Collision with player
		if g.Player.isColliding(potion.Sprite) {
			g.Player.Health += potion.AmtHeal
			fmt.Println("Player Healed:", g.Player.Health)
			// remove potion from potions
			g.potions = append(g.potions[:i], g.potions[i+1:]...)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 120, G: 180, B: 255, A: 255})
	opts := ebiten.DrawImageOptions{}

	// loop over layers
	for _, layer := range g.TilemapJSON.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(g.TilemapImage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image), &opts)
			opts.GeoM.Reset()
		}
	}

	// set translation to players position
	opts.GeoM.Translate(g.Player.X, g.Player.Y)
	screen.DrawImage(
		g.Player.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts)

	opts.GeoM.Reset()

	for _, sprite := range g.sprites {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}
	opts.GeoM.Reset()

	opts.GeoM.Reset()

	for _, sprite := range g.potions {
		opts.GeoM.Translate(sprite.X, sprite.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}
	opts.GeoM.Reset()
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	playerImg, _, err := ebitenutil.NewImageFromFile("./assets/images/ninja.png")
	if err != nil {
		log.Fatal(err)
	}
	skeleton, _, err := ebitenutil.NewImageFromFile("./assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}
	potionImg, _, err := ebitenutil.NewImageFromFile("./assets/images/health.png")
	if err != nil {
		log.Fatal(err)
	}
	tileMapImage, _, err := ebitenutil.NewImageFromFile("./assets/images/TileSetFloor.png")
	player := &Player{
		Sprite: &Sprite{
			Img: playerImg,
			X:   200,
			Y:   100,
		},
		Health: 3,
	}

	tileMapJson, err := NewTilemapLayerJSON("./assets/maps/spawn.json")

	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		Player: player,
		potions: []*Potion{
			{
				&Sprite{
					Img: potionImg,
					X:   150,
					Y:   100,
				},
				1.0,
			}, {
				&Sprite{
					Img: potionImg,
					X:   180,
					Y:   10,
				},
				1.0,
			},
		},
		sprites: []*Enemy{
			{
				&Sprite{
					Img: skeleton,
					X:   100,
					Y:   100,
				},
				true,
			},
			{
				&Sprite{
					Img: skeleton,
					X:   60,
					Y:   70,
				},
				false,
			},
			{
				Sprite: &Sprite{
					Img: skeleton,
					X:   120,
					Y:   180,
				},
				FollowsPlayer: false,
			},
		},
		TilemapJSON:  tileMapJson,
		TilemapImage: tileMapImage,
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
