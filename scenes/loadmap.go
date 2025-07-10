package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"rpg-go/collisions"
	"rpg-go/components"
	"rpg-go/constants"
	"rpg-go/entities"
	"rpg-go/tilemap"
)

// LoadMap limpa o estado do mapa antigo e carrega um novo.
func (g *GameScene) LoadMap(mapPath string, targetSpawn string) {
	// Limpa entidades e colisões do mapa anterior
	g.enemies = make([]*entities.Enemy, 0)
	g.potions = make([]*entities.Potion, 0)
	g.dummies = make([]*entities.TrainingDummy, 0)

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

	mapWidthPixels := g.TilemapJSON.Layers[0].Width * constants.Tilesize
	mapHeightPixels := g.TilemapJSON.Layers[0].Height * constants.Tilesize
	g.CollisionGrid = collisions.NewGrid(mapWidthPixels, mapHeightPixels)

	colliderCount := 0

	spawnPoints := make(map[string]image.Point)

	for _, layer := range g.TilemapJSON.Layers {
		if layer.Type == "objectgroup" {
			log.Printf("Processando camada de objetos: '%s'", layer.Name)
			if layer.Name == "collisions" {
				for _, col := range layer.Objects {
					x := col.X
					y := col.Y
					w := col.Width
					h := col.Height

					if col.GID > 0 {
						y = y - h
					}

					collider := image.Rect(int(x), int(y), int(x+w), int(y+h))
					g.CollisionGrid.Insert(&collider)
					colliderCount++
				}
			}
			for _, obj := range layer.Objects {
				switch obj.Type {
				case "player_spawn":
					spawnName := obj.Name
					if spawnName == "" {
						spawnName = "default"
					}
					spawnPoints[spawnName] = image.Point{X: int(obj.X), Y: int(obj.Y)}

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
						Sprite: &entities.Sprite{
							Img: g.assets.SkeletonImg,
							X:   obj.X,
							Y:   obj.Y,
						},
						FollowsPlayer: follows,
						CombatComp:    components.NewEnemieCombat(3, 1, 60), // Cooldown de 1s (60 ticks)
					}
					g.enemies = append(g.enemies, newEnemy)

				case "training_dummy":

					newDummy := entities.NewTrainingDummy(obj.X, obj.Y, g.assets.DummyImg)
					g.dummies = append(g.dummies, newDummy)

				case "potion_spawn":

					amount, found := tilemap.GetIntProperty("amount", obj.Properties)
					if !found {
						amount = 2
					}
					newPotion := &entities.Potion{
						Sprite: &entities.Sprite{
							Img: g.assets.PotionImg,
							X:   obj.X,
							Y:   obj.Y,
						},
						AmtHeal: uint(amount),
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

func (g *GameScene) debugDrawColliders(screen *ebiten.Image) {
	sw, sh := screen.Size()
	camRect := image.Rect(
		int(-g.Camera.X), int(-g.Camera.Y),
		int(-g.Camera.X)+sw, int(-g.Camera.Y)+sh,
	)

	// Pega os colisores que estão na área visível da câmera
	nearbyColliders := g.CollisionGrid.GetNearbyColliders(camRect)

	// Log de depuração para a função Draw
	if len(nearbyColliders) > 0 && ebiten.IsKeyPressed(ebiten.KeyC) { // Pressione C para ver o log
		log.Printf("Debug Draw: Encontrados %d colisores próximos para desenhar.", len(nearbyColliders))
	}

	for _, collider := range nearbyColliders {
		camX, camY := g.Camera.X, g.Camera.Y
		x := float32(collider.Min.X) + float32(camX)
		y := float32(collider.Min.Y) + float32(camY)
		w := float32(collider.Dx())
		h := float32(collider.Dy())
		vector.StrokeRect(screen, x, y, w, h, 1, color.RGBA{R: 255, G: 0, B: 0, A: 255}, false)
	}
}
