package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"log"
	"math"
	"rpg-go/entities"
	"sort"
)

type Game struct {
	Player       *entities.Player
	sprites      []*entities.Enemy
	potions      []*entities.Potion
	TilemapJSON  *TilemapJSON
	TilemapImage *ebiten.Image
	Camera       *Camera
}

func (g *Game) Update() error {
	// React to keyPresses

	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx += 2
	}
	// Normalize movement
	// Magnitude
	magnitude := math.Sqrt(dx*dx + dy*dy)

	if magnitude > 0.0 {
		speed := 2.0
		dx = (dx / magnitude) * speed
		dy = (dy / magnitude) * speed
	}

	g.Player.X += dx
	g.Player.Y += dy

	for _, sprite := range g.sprites {
		if sprite.FollowsPlayer {
			if sprite.X < math.Floor(g.Player.X) {
				sprite.X += 1
			} else if sprite.X > math.Floor(g.Player.X) {
				sprite.X -= 1
			}
			if sprite.Y < math.Floor(g.Player.Y) {
				sprite.Y += 1
			} else if sprite.Y > math.Floor(g.Player.Y) {
				sprite.Y -= 1
			}
		}
	}

	for i, potion := range g.potions {
		// Collision with player
		if g.Player.IsColliding(potion.Sprite) {
			g.Player.Health += potion.AmtHeal
			fmt.Println("Player Healed:", g.Player.Health)
			// remove potion from potions
			g.potions = append(g.potions[:i], g.potions[i+1:]...)
		}
	}

	g.Camera.FollowTarget(g.Player.X+8, g.Player.Y+8, 320, 240)
	g.Camera.Constrain(float64(g.TilemapJSON.Layers[0].Width*16), float64(g.TilemapJSON.Layers[0].Height)*16.0, 320, 240)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 120, G: 180, B: 255, A: 255})
	opts := ebiten.DrawImageOptions{}

	// loop over layers
	for _, layer := range g.TilemapJSON.Layers {
		// loop over the tiles
		for index, id := range layer.Data {

			// get the position of pixel
			x := index % layer.Width
			y := index / layer.Width

			//convert to pixel position
			x *= 16
			y *= 16

			// get the position of the image where the tile id is
			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			// convert to pixel position
			srcX *= 16
			srcY *= 16

			// set the draw image options to draw on x,y
			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(g.Camera.X, g.Camera.Y)

			//draw the image
			screen.DrawImage(
				// crop the image
				g.TilemapImage.SubImage(
					image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image), &opts)
			// reset the draw
			opts.GeoM.Reset()
		}
	}

	// set translation to players position
	opts.GeoM.Translate(g.Player.X, g.Player.Y)

	ent := []*entities.Sprite{g.Player.Sprite}
	for _, sprite := range g.sprites {
		ent = append(ent, sprite.Sprite)
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

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
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

	tileMapJson, err := NewTilemapLayerJSON("./assets/maps/spawn.json")

	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		Player: player,
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
		sprites: []*entities.Enemy{
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
		TilemapJSON:  tileMapJson,
		TilemapImage: tileMapImage,
		Camera:       NewCamera(0, 0),
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
