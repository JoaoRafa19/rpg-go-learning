package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"rpg-go/animations"
	"rpg-go/components"
	"rpg-go/constants"
	"rpg-go/entities"
	"rpg-go/spritesheet"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	animationFram     int
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	TilemapJSON       *TilemapJSON
	Tilesets          []Tileset
	TilemapImage      *ebiten.Image
	Camera            *Camera
	Coliders          []Collisor
}

func NewGame() *Game {
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
		Animations: map[entities.PlayerState]*animations.Animation{
			entities.Up:    animations.NewAnimation(5, 13, 4, 20.0),
			entities.Down:  animations.NewAnimation(4, 12, 4, 20),
			entities.Left:  animations.NewAnimation(6, 14, 4, 20),
			entities.Right: animations.NewAnimation(7, 15, 4, 20),
		},
		CombatComp: components.NewBasicCombat(3, 1),
		Sprite: &entities.Sprite{
			Img: playerImg,
			X:   200,
			Y:   100,
		},
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}
	playerSpritesheet := spritesheet.NewSpriteSheet(4, 7, constants.Tilesize)

	colliders := tilemapJSON.GenCollisors()
	return &Game{
		Tilesets:          tilesets,
		player:            player,
		Coliders:          colliders,
		playerSpriteSheet: playerSpritesheet,

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
				Sprite: &entities.Sprite{
					Img: skeleton,
					X:   100,
					Y:   200,
				},
				CombatComp:    components.NewEnemieCombat(3, 1, 30),
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeleton,
					X:   60,
					Y:   70,
				},
				FollowsPlayer: false,
				CombatComp:    components.NewEnemieCombat(3, 1, 30),
			},
			{
				Sprite: &entities.Sprite{
					Img: skeleton,
					X:   120,
					Y:   180,
				},
				FollowsPlayer: false,
				CombatComp:    components.NewEnemieCombat(3, 1, 30),
			},
		},
		TilemapJSON:  tilemapJSON,
		TilemapImage: tileMapImage,
		Camera:       NewCamera(0, 0),
	}

}

func (g *Game) Update() error {

	// React to keyPresses
	g.player.Dx = 0.0
	g.player.Dy = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Dy = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Dx += 2
	}
	// Normalize movement
	// Magnitude
	magnitude := math.Sqrt(g.player.Dx*g.player.Dx + g.player.Dy*g.player.Dy)

	if magnitude > 0.0 {
		speed := 2.0
		g.player.Dx = (g.player.Dx / magnitude) * speed
		g.player.Dy = (g.player.Dy / magnitude) * speed
	}

	g.player.X += g.player.Dx
	CheckCollisionsHorizontaly(g.player.Sprite, g.Coliders)

	g.player.Y += g.player.Dy
	CheckCollisionsVerticaly(g.player.Sprite, g.Coliders)

	activeAnimation := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnimation != nil {
		activeAnimation.Update()
	}

	for _, sprite := range g.enemies {
		sprite.Dx = 0.0
		sprite.Dy = 0.0
		if sprite.FollowsPlayer {
			if sprite.X < math.Floor(g.player.X) {
				sprite.Dx += 1
			} else if sprite.X > math.Floor(g.player.X) {
				sprite.Dx -= 1
			}
			if sprite.Y < math.Floor(g.player.Y) {
				sprite.Dy += 1
			} else if sprite.Y > math.Floor(g.player.Y) {
				sprite.Dy -= 1
			}
		}
		sprite.X += sprite.Dx
		CheckCollisionsHorizontaly(sprite.Sprite, g.Coliders)
		sprite.Y += sprite.Dy
		CheckCollisionsVerticaly(sprite.Sprite, g.Coliders)
	}

	cX, cY := ebiten.CursorPosition()
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)

	// Converter para float64 para cálculos precisos

	worldX := float64(cX) - g.Camera.X
	worldY := float64(cY) - g.Camera.Y
	deadEnemies := make(map[int]struct{})

	g.player.CombatComp.Update()
	pRect := image.Rect(
		int(g.player.X),
		int(g.player.Y),
		int(g.player.X)+constants.Tilesize,
		int(g.player.Y)+constants.Tilesize,
	)

	for idx, enemy := range g.enemies {
		enemy.CombatComp.Update()
		rect := image.Rect(
			int(enemy.X),
			int(enemy.Y),
			int(enemy.X)+constants.Tilesize,
			int(enemy.Y)+constants.Tilesize,
		)

		if rect.Overlaps(pRect) {
			if enemy.CombatComp.Attack() {
				g.player.CombatComp.Damage(enemy.CombatComp.AttackPower())
				if g.player.CombatComp.Health() <= 0 {
					//fmt.Println("player died!!!")
				}
			}
		}

		if worldX > float64(rect.Min.X) && worldX < float64(rect.Max.X) && worldY > float64(rect.Min.Y) && worldY < float64(rect.Max.Y) {

			if clicked && math.Sqrt(math.Pow(float64(worldX)-g.player.X+(constants.Tilesize/2), 2)+math.Pow(float64(worldY)-g.player.Y+(constants.Tilesize/2), 2)) < constants.Tilesize*5 {
				enemy.CombatComp.Damage(g.player.CombatComp.AttackPower())

				if enemy.CombatComp.Health() <= 0 {
					deadEnemies[idx] = struct{}{}
				}
			}
		}
	}

	if len(deadEnemies) > 0 {
		newEnemies := make([]*entities.Enemy, 0)
		for index, enemy := range g.enemies {
			if _, exts := deadEnemies[index]; !exts {
				newEnemies = append(newEnemies, enemy)
			}
		}
		g.enemies = newEnemies
	}

	g.Camera.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	g.Camera.Constrain(float64(g.TilemapJSON.Layers[0].Width*constants.Tilesize), float64(g.TilemapJSON.Layers[0].Height)*constants.Tilesize, 320, 240)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 120, G: 180, B: 255, A: 255})
	opts := ebiten.DrawImageOptions{}

	// loop over layers
	for layerIndex, layer := range g.TilemapJSON.Layers {
		if len(layer.Objects) > 0 {
			for _, colisor := range layer.Objects {
				vector.StrokeRect(
					screen,
					float32(colisor.X)+float32(g.Camera.X),
					float32(colisor.Y)+float32(g.Camera.Y),
					float32(colisor.Width),
					float32(colisor.Width),
					1.0,
					color.RGBA{255, 0, 0, 255},
					true,
				)
			}
		}
		// loop over the tiles
		for index, id := range layer.Data {

			if id == 0 {
				continue
			}

			// get the position of pixel
			x := index % layer.Width
			y := index / layer.Width

			//convert to pixel position
			x *= constants.Tilesize
			y *= constants.Tilesize

			img := g.Tilesets[layerIndex].Img(id)
			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + constants.Tilesize))
			opts.GeoM.Translate(g.Camera.X, g.Camera.Y)

			screen.DrawImage(img, &opts)

			opts.GeoM.Reset()
		}
	}

	// set translation to players position

	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.Camera.X, g.Camera.Y)

	playerFrame := 0
	activeAnimation := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnimation != nil {
		playerFrame = activeAnimation.Frame()
	}
	screen.DrawImage(
		g.player.Img.SubImage(
			g.playerSpriteSheet.Rect(playerFrame),
		).(*ebiten.Image),
		&opts,
	)
	opts.GeoM.Reset()
	ent := []*entities.Sprite{}
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
				image.Rect(0, 0, constants.Tilesize, constants.Tilesize),
			).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}

	//for _, collider := range g.TilemapJSON.Tilesets. {

	//}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}
