package scenes

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"rpg-go/animations"
	"rpg-go/camera"
	"rpg-go/components"
	"rpg-go/constants"
	"rpg-go/entities"
	"rpg-go/spritesheet"
	"rpg-go/tilemap"
	"rpg-go/tileset"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameScene struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	animationFram     int
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	TilemapJSON       *tilemap.TilemapJSON
	Tilesets          []tileset.Tileset
	TilemapImage      *ebiten.Image
	Camera            *camera.Camera
	Colliders         []image.Rectangle
	loaded            bool
}

// IsLoaded implements Scene.
func (g *GameScene) IsLoaded() bool {
	return g.loaded
}

func NewGameScene() *GameScene {
	return &GameScene{
		player:            nil,
		playerSpriteSheet: nil,
		enemies:           make([]*entities.Enemy, 0),
		potions:           make([]*entities.Potion, 0),
		TilemapJSON:       nil,
		Tilesets:          nil,
		TilemapImage:      nil,
		Camera:            nil,
		Colliders:         nil,
		loaded:            false,
	}
}

// Draw implements Scene.
func (g *GameScene) Draw(screen *ebiten.Image) {
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

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

// FirstLoad implements Scene.
func (g *GameScene) FirstLoad() {
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
	g.player = &entities.Player{
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
	g.TilemapImage = tileMapImage

	tilemapJSON, err := tilemap.NewTilemapJSON("assets/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}
	g.TilemapJSON = tilemapJSON

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}
	g.Tilesets = tilesets

	g.playerSpriteSheet = spritesheet.NewSpriteSheet(4, 7, constants.Tilesize)

	g.Colliders = []image.Rectangle{
		image.Rect(100, 100, 116, 116),
		image.Rect(10, 10, 116, 116),
	}

	g.potions = []*entities.Potion{
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
	}

	g.enemies = []*entities.Enemy{
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
	}

	g.Camera = camera.NewCamera(0, 0)
	g.loaded = true
}

// OnEnter implements Scene.
func (g *GameScene) OnEnter() {

}

// OnExit implements Scene.
func (g *GameScene) OnExit() {

}

func CheckCollisionsVerticaly(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+constants.Tilesize, int(sprite.Y)+constants.Tilesize)) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - constants.Tilesize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}
func CheckCollisionsHorizontaly(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+constants.Tilesize, int(sprite.Y)+constants.Tilesize)) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - constants.Tilesize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

// Update implements Scene.
func (g *GameScene) Update() SceneId {

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return PauseSceneId
	}
	

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
	CheckCollisionsHorizontaly(g.player.Sprite, g.Colliders)

	g.player.Y += g.player.Dy
	CheckCollisionsVerticaly(g.player.Sprite, g.Colliders)

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
		CheckCollisionsHorizontaly(sprite.Sprite, g.Colliders)
		sprite.Y += sprite.Dy
		CheckCollisionsVerticaly(sprite.Sprite, g.Colliders)
	}

	cX, cY := ebiten.CursorPosition()
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)

	// Converter para float64 para cálculos precisos

	worldX := float64(cX) - g.Camera.X
	worldY := float64(cY) - g.Camera.Y
	deadEnemies := make(map[int]struct{})
	g.Colliders = []image.Rectangle{}
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
				fmt.Println(fmt.Sprintf("DAMAGE::%d::%d", g.player.CombatComp.Health()))
				if g.player.CombatComp.Health() <= 0 {
					fmt.Println("player died!!!")
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
	return GameSceneId
}

var _ Scene = (*GameScene)(nil)
