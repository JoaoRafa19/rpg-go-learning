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
	"rpg-go/hud"
	"rpg-go/spritesheet"
	"rpg-go/tilemap"
	"rpg-go/tileset"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameScene struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	TilemapJSON       *tilemap.TilemapJSON
	Tilesets          []*tileset.Tileset
	Camera            *camera.Camera
	Colliders         []image.Rectangle
	loaded            bool
	hud               *hud.HUD

	// Assets que não mudam entre mapas
	playerImg   *ebiten.Image
	skeletonImg *ebiten.Image
	potionImg   *ebiten.Image
}

func NewGameScene() *GameScene {
	return &GameScene{
		enemies:   make([]*entities.Enemy, 0),
		potions:   make([]*entities.Potion, 0),
		Colliders: make([]image.Rectangle, 0),
		loaded:    false,
	}
}

func (g *GameScene) FirstLoad() {
	var err error
	g.playerImg, _, err = ebitenutil.NewImageFromFile("./assets/images/ninja.png")
	if err != nil {
		log.Fatal(err)
	}
	g.skeletonImg, _, err = ebitenutil.NewImageFromFile("./assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}
	g.potionImg, _, err = ebitenutil.NewImageFromFile("./assets/images/health.png")
	if err != nil {
		log.Fatal(err)
	}

	g.player = &entities.Player{
		Animations: map[entities.PlayerState]*animations.Animation{
			entities.Up:    animations.NewAnimation(5, 13, 4, 20.0),
			entities.Down:  animations.NewAnimation(4, 12, 4, 20),
			entities.Left:  animations.NewAnimation(6, 14, 4, 20),
			entities.Right: animations.NewAnimation(7, 15, 4, 20),
		},
		CombatComp: components.NewBasicCombat(10, 1), // Aumentei a vida para 10
		Sprite: &entities.Sprite{
			Img: g.playerImg,
		},
	}

	g.playerSpriteSheet = spritesheet.NewSpriteSheet(4, 7, constants.Tilesize)
	g.hud = hud.NewHUD(g.player.CombatComp)
	g.Camera = camera.NewCamera(0, 0)

	// Carrega o mapa inicial e posiciona o jogador
	g.LoadMap("assets/maps/spawn.json", "default")

	g.loaded = true
}

// LoadMap limpa o estado do mapa antigo e carrega um novo.
func (g *GameScene) LoadMap(mapPath string, targetSpawn string) {
	// Limpa entidades e colisões do mapa anterior
	g.enemies = make([]*entities.Enemy, 0)
	g.potions = make([]*entities.Potion, 0)
	g.Colliders = make([]image.Rectangle, 0)

	// Carrega o JSON do novo mapa
	tilemapJSON, err := tilemap.NewTilemapJSON(mapPath)
	if err != nil {
		log.Fatalf("Could not load map %s: %v", mapPath, err)
	}
	g.TilemapJSON = tilemapJSON

	// Gera os tilesets para o novo mapa
	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}
	g.Tilesets = tilesets

	spawnPoints := make(map[string]image.Point)

	// Carrega os objetos do mapa (inimigos, poções, colisões, spawns)
	for _, layer := range g.TilemapJSON.Layers {
		if layer.Type == "objectgroup" {
			for _, obj := range layer.Objects {
				switch obj.Type {
				case "player_spawn":
					spawnName := obj.Name
					if spawnName == "" {
						spawnName = "default"
					}
					spawnPoints[spawnName] = image.Point{X: int(obj.X), Y: int(obj.Y)}

				case "collision":
					g.Colliders = append(g.Colliders, image.Rect(int(obj.X), int(obj.Y), int(obj.X+obj.Width), int(obj.Y+obj.Height)))

				case "enemy_spawn":
					follows := false
					for _, prop := range obj.Properties {
						if prop.Name == "followsPlayer" {
							if val, ok := prop.Value.(bool); ok {
								follows = val
							}
						}
					}
					newEnemy := &entities.Enemy{
						Sprite:        &entities.Sprite{Img: g.skeletonImg, X: obj.X, Y: obj.Y},
						FollowsPlayer: follows,
						CombatComp:    components.NewEnemieCombat(3, 1, 60), // Cooldown de 1s (60 ticks)
					}
					g.enemies = append(g.enemies, newEnemy)

				case "potion_spawn":
					newPotion := &entities.Potion{
						Sprite:  &entities.Sprite{Img: g.potionImg, X: obj.X, Y: obj.Y},
						AmtHeal: 3, // Cura 3 de vida
					}
					g.potions = append(g.potions, newPotion)
				}
			}
		}
	}

	// Posiciona o jogador no ponto de spawn correto
	spawnPos, found := spawnPoints[targetSpawn]
	if !found {
		// Se o spawn alvo não for encontrado, usa o "default"
		spawnPos = spawnPoints["default"]
	}
	g.player.X = float64(spawnPos.X)
	g.player.Y = float64(spawnPos.Y)
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{144, 208, 128, 255}) // Um verde mais agradável
	opts := &ebiten.DrawImageOptions{}

	// --- Desenha o Mapa ---
	camX, camY := g.Camera.X, g.Camera.Y
	for _, layer := range g.TilemapJSON.Layers {
		if layer.Type != "tilelayer" {
			continue
		}
		for i, tileID := range layer.Data {
			if tileID == 0 {
				continue
			}

			tileset := g.findTilesetForTile(tileID)
			if tileset == nil {
				continue
			}

			x := float64((i % layer.Width) * constants.Tilesize)
			y := float64((i / layer.Width) * constants.Tilesize)

			opts.GeoM.Reset()
			opts.GeoM.Translate(x, y)
			opts.GeoM.Translate(camX, camY)

			screen.DrawImage(tileset.Img(int(tileID)), opts)
		}
	}

	// --- Desenha Entidades com Ordenação Y (Profundidade) ---
	var drawables []entities.Drawable
	drawables = append(drawables, g.player)
	for _, e := range g.enemies {
		drawables = append(drawables, e)
	}
	for _, p := range g.potions {
		drawables = append(drawables, p)
	}

	sort.Slice(drawables, func(i, j int) bool {
		return drawables[i].GetY() < drawables[j].GetY()
	})

	for _, d := range drawables {
		d.Draw(screen, g.Camera, g.playerSpriteSheet)
	}

	// --- Desenha o HUD ---
	g.hud.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

func (g *GameScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return PauseSceneId
	}

	// 1. Lidar com a entrada e movimento do jogador
	g.handlePlayerMovement()

	// 2. Atualizar animações
	activeAnimation := g.player.ActiveAnimation(int(g.player.Dx), int(g.player.Dy))
	if activeAnimation != nil {
		activeAnimation.Update()
	}

	// 3. Atualizar inimigos
	g.updateEnemies()

	// 4. Lidar com combate
	g.handleCombat()

	// 5. Lidar com itens coletáveis
	g.handleCollectibles()

	// 6. Checar transições de mapa
	if nextMap, nextSpawn := g.checkMapTransitions(); nextMap != "" {
		g.LoadMap(nextMap, nextSpawn)
		// Retornar aqui para o próximo frame começar com o mapa já carregado
		return GameSceneId
	}

	// 7. Atualizar a câmera
	g.Camera.FollowTarget(g.player.X+8, g.player.Y+8, 320, 240)
	mapWidth := float64(g.TilemapJSON.Layers[0].Width * constants.Tilesize)
	mapHeight := float64(g.TilemapJSON.Layers[0].Height * constants.Tilesize)
	g.Camera.Constrain(mapWidth, mapHeight, 320, 240)

	return GameSceneId
}

func (g *GameScene) handlePlayerMovement() {
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

	if g.player.Dx != 0 || g.player.Dy != 0 {
		magnitude := math.Sqrt(g.player.Dx*g.player.Dx + g.player.Dy*g.player.Dy)
		speed := 2.0
		g.player.Dx = (g.player.Dx / magnitude) * speed
		g.player.Dy = (g.player.Dy / magnitude) * speed
	}

	g.player.X += g.player.Dx

	g.player.Y += g.player.Dy
}

func (g *GameScene) updateEnemies() {
	for _, enemy := range g.enemies {
		enemy.CombatComp.Update()

		enemy.Dx = 0.0
		enemy.Dy = 0.0
		if enemy.FollowsPlayer {
			// Lógica de perseguição simples
			if math.Abs(enemy.X-g.player.X) > 1 {
				if enemy.X < g.player.X {
					enemy.Dx = 1
				} else {
					enemy.Dx = -1
				}
			}
			if math.Abs(enemy.Y-g.player.Y) > 1 {
				if enemy.Y < g.player.Y {
					enemy.Dy = 1
				} else {
					enemy.Dy = -1
				}
			}
		}

		enemy.X += enemy.Dx
		enemy.Y += enemy.Dy
	}
}

func (g *GameScene) handleCombat() {
	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	g.player.CombatComp.Update()

	deadEnemies := map[int]struct{}{}
	pRect := image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X)+constants.Tilesize, int(g.player.Y)+constants.Tilesize)

	for idx, enemy := range g.enemies {
		enemyRect := image.Rect(int(enemy.X), int(enemy.Y), int(enemy.X)+constants.Tilesize, int(enemy.Y)+constants.Tilesize)

		// Combate: Inimigo ataca o Jogador
		if enemyRect.Overlaps(pRect) {
			if enemy.CombatComp.Attack() {
				g.player.CombatComp.Damage(enemy.CombatComp.AttackPower())
				fmt.Printf("Player took damage! Health: %d\n", g.player.CombatComp.Health())
				if g.player.CombatComp.Health() <= 0 {
					fmt.Println("PLAYER DIED! Game Over.")
					// TODO: Implementar lógica de Game Over (ex: ir para uma cena de Game Over)
				}
			}
		}

		// Combate: Jogador ataca o Inimigo
		if clicked {
			cX, cY := ebiten.CursorPosition()
			worldX, worldY := float64(cX)-g.Camera.X, float64(cY)-g.Camera.Y

			// Verifica se o clique foi no inimigo
			if worldX >= enemy.X && worldX < enemy.X+constants.Tilesize && worldY >= enemy.Y && worldY < enemy.Y+constants.Tilesize {
				// Verifica o alcance do ataque
				distance := math.Sqrt(math.Pow(enemy.X-g.player.X, 2) + math.Pow(enemy.Y-g.player.Y, 2))
				if distance < float64(constants.Tilesize)*2.5 { // Alcance de ataque de 2.5 tiles
					enemy.CombatComp.Damage(g.player.CombatComp.AttackPower())
					fmt.Printf("Enemy took damage! Health: %d\n", enemy.CombatComp.Health())
					if enemy.CombatComp.Health() <= 0 {
						deadEnemies[idx] = struct{}{}
					}
				}
			}
		}
	}

	// Remove inimigos mortos
	if len(deadEnemies) > 0 {
		newEnemies := make([]*entities.Enemy, 0, len(g.enemies)-len(deadEnemies))
		for i, enemy := range g.enemies {
			if _, isDead := deadEnemies[i]; !isDead {
				newEnemies = append(newEnemies, enemy)
			}
		}
		g.enemies = newEnemies
	}
}

func (g *GameScene) handleCollectibles() {
	potionsToCollect := []int{}
	pRect := image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X)+constants.Tilesize, int(g.player.Y)+constants.Tilesize)

	for i, potion := range g.potions {
		potionRect := image.Rect(int(potion.X), int(potion.Y), int(potion.X)+constants.Tilesize, int(potion.Y)+constants.Tilesize)
		if pRect.Overlaps(potionRect) {
			if g.player.CombatComp.Health() < g.player.CombatComp.MaxHealth() {
				g.player.CombatComp.Heal(int(potion.AmtHeal))
				fmt.Printf("Player healed! Current Health: %d\n", g.player.CombatComp.Health())
				potionsToCollect = append(potionsToCollect, i)
			}
		}
	}

	if len(potionsToCollect) > 0 {
		sort.Sort(sort.Reverse(sort.IntSlice(potionsToCollect)))
		for _, index := range potionsToCollect {
			g.potions = append(g.potions[:index], g.potions[index+1:]...)
		}
	}
}

func (g *GameScene) checkMapTransitions() (string, string) {
	pRect := image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X)+constants.Tilesize, int(g.player.Y)+constants.Tilesize)

	for _, layer := range g.TilemapJSON.Layers {
		if layer.Type == "objectgroup" {
			for _, obj := range layer.Objects {
				if obj.Type == "transition" {
					objRect := image.Rect(int(obj.X), int(obj.Y), int(obj.X+obj.Width), int(obj.Y+obj.Height))
					if pRect.Overlaps(objRect) {
						targetMap := ""
						targetSpawn := "default" // Padrão
						for _, prop := range obj.Properties {
							if prop.Name == "targetMap" {
								targetMap = prop.Value.(string)
							}
							if prop.Name == "targetSpawn" {
								targetSpawn = prop.Value.(string)
							}
						}

						if targetMap != "" {
							return "assets/maps/" + targetMap, targetSpawn
						}
					}
				}
			}
		}
	}
	return "", ""
}

func (g *GameScene) IsLoaded() bool {
	return g.loaded
}

// Helper para encontrar o tileset correto para um determinado ID de tile
func (g *GameScene) findTilesetForTile(tileID int) *tileset.Tileset {
	for i := len(g.Tilesets) - 1; i >= 0; i-- {
		if tileID >= g.Tilesets[i].FirstGid {
			return g.Tilesets[i]
		}
	}
	return nil
}

// ... (as funções IsLoaded, OnEnter, OnExit, CheckCollisions permanecem as mesmas)
func (g *GameScene) OnEnter() {}
func (g *GameScene) OnExit()  {}

// ... (código das funções de checagem de colisão) ...
// (Este código pode permanecer igual ao seu)

// É necessário adicionar interfaces e métodos nas suas entidades para o sorteio Y e desenho
// Exemplo em entities/sprite.go ou em cada entidade:
// type Drawable interface {
//     GetY() float64
//     Draw(screen *ebiten.Image, camera *camera.Camera, sheet *spritesheet.SpriteSheet)
// }
