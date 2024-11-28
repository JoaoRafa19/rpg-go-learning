package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"rpg-go/entities"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Player       *entities.Player
	enemies      []*entities.Enemy
	potions      []*entities.Potion
	TilemapJSON  *TilemapJSON
	Tilesets     []Tileset
	TilemapImage *ebiten.Image
	Camera       *Camera
	Colliders    []image.Rectangle
}

func CheckCollisionsVerticaly(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+16, int(sprite.Y)+16)) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - 16.0
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}
func CheckCollisionsHorizontaly(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+16, int(sprite.Y)+16)) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - 16.0
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func (g *Game) Update() error {
	// React to keyPresses
	g.Player.Dx = 0.0
	g.Player.Dy = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Player.Dy = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Player.Dx += 2
	}
	// Normalize movement
	// Magnitude
	magnitude := math.Sqrt(g.Player.Dx*g.Player.Dx + g.Player.Dy*g.Player.Dy)

	if magnitude > 0.0 {
		speed := 2.0
		g.Player.Dx = (g.Player.Dx / magnitude) * speed
		g.Player.Dy = (g.Player.Dy / magnitude) * speed
	}

	g.Player.X += g.Player.Dx
	CheckCollisionsHorizontaly(g.Player.Sprite, g.Colliders)
	
	g.Player.Y += g.Player.Dy
	CheckCollisionsVerticaly(g.Player.Sprite, g.Colliders)

	for _, sprite := range g.enemies {
		sprite.Dx = 0.0
		sprite.Dy = 0.0
		if sprite.FollowsPlayer {
			if sprite.X < math.Floor(g.Player.X) {
				sprite.Dx += 1
			} else if sprite.X > math.Floor(g.Player.X) {
				sprite.Dx -= 1
			}
			if sprite.Y < math.Floor(g.Player.Y) {
				sprite.Dy += 1
			} else if sprite.Y > math.Floor(g.Player.Y) {
				sprite.Dy -= 1
			}
		}
		sprite.X += sprite.Dx
		CheckCollisionsHorizontaly(sprite.Sprite, g.Colliders)
		sprite.Y += sprite.Dy
		CheckCollisionsVerticaly(sprite.Sprite, g.Colliders)
	}

	g.Camera.FollowTarget(g.Player.X+8, g.Player.Y+8, 320, 240)
	g.Camera.Constrain(float64(g.TilemapJSON.Layers[0].Width*16), float64(g.TilemapJSON.Layers[0].Height)*16.0, 320, 240)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 120, G: 180, B: 255, A: 255})
	opts := ebiten.DrawImageOptions{}

	// loop over layers
	for layerIndex, layer := range g.TilemapJSON.Layers {
		// loop over the tiles
		for index, id := range layer.Data {

			if id == 0 {
				continue
			}

			// get the position of pixel
			x := index % layer.Width
			y := index / layer.Width

			//convert to pixel position
			x *= 16
			y *= 16

			img := g.Tilesets[layerIndex].Img(id)
			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))
			opts.GeoM.Translate(g.Camera.X, g.Camera.Y)

			screen.DrawImage(img, &opts)

			opts.GeoM.Reset()
		}
	}

	// set translation to players position

	ent := []*entities.Sprite{g.Player.Sprite}
	for _, enemy := range g.enemies {
		ent = append(ent, enemy.Sprite)
	}
	for _, potion := range g.potions {
		ent = append(ent, potion.Sprite)
	}

	// sort by Y
	sort.Slice(ent, func(i, j int) bool {
		return ent[i].Y < ent[j].Y
	})

	for _, sprite := range ent {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.Camera.X, g.Camera.Y)

		screen.DrawImage(
			sprite.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

	for _, collider := range g.Colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.Camera.X),
			float32(collider.Min.Y)+float32(g.Camera.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			true,
		)
	}

	/*
		screen.DrawImage(
			g.Player.Img.SubImage(
				image.Rect(0, 0, 16, 16),
			).(*ebiten.Image),
			&opts)

		opts.GeoM.Reset()

		for _, sprite := range g.sprites {
			opts.GeoM.Translate(sprite.X, sprite.Y)
			opts.GeoM.Translate(g.Camera.X, g.Camera.Y)

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
			opts.GeoM.Translate(g.Camera.X, g.Camera.Y)
			screen.DrawImage(
				sprite.Img.SubImage(
					image.Rect(0, 0, 16, 16),
				).(*ebiten.Image),
				&opts,
			)
			opts.GeoM.Reset()
		}
		opts.GeoM.Reset()

	*/

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("RPG Go!")
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
	player := &entities.Player{
		Sprite: &entities.Sprite{
			Img: playerImg,
			X:   200,
			Y:   100,
		},
		Health: 3,
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		Tilesets: tilesets,
		Player:   player,
		Colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   150,
					Y:   100,
				},
				AmtHeal: 1.0,
			},
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   180,
					Y:   10,
				},
				AmtHeal: 1.0,
			},
		},
		enemies: []*entities.Enemy{
			{
				&entities.Sprite{
					Img: skeleton,
					X:   100,
					Y:   100,
				},
				true,
			},
			{
				&entities.Sprite{
					Img: skeleton,
					X:   60,
					Y:   70,
				},
				false,
			},
			{
				&entities.Sprite{
					Img: skeleton,
					X:   120,
					Y:   180,
				},
				false,
			},
		},
		TilemapJSON:  tilemapJSON,
		TilemapImage: tileMapImage,
		Camera:       NewCamera(0, 0),
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
