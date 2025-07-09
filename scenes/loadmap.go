package scenes

import (
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"rpg-go/components"
	"rpg-go/entities"
	"rpg-go/tilemap"
)

// LoadMap limpa o estado do mapa antigo e carrega um novo.
func (g *GameScene) LoadMap(mapPath string, targetSpawn string) {
	// Limpa entidades e colisões do mapa anterior
	g.enemies = make([]*entities.Enemy, 0)
	g.potions = make([]*entities.Potion, 0)
	g.Colliders = make([]image.Rectangle, 0)
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

	spawnPoints := make(map[string]image.Point)

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
					skeletonImg, _, err := ebitenutil.NewImageFromFile("./assets/images/skeleton.png")
					if err != nil {
						log.Fatal(err)
					}
					newEnemy := &entities.Enemy{
						Sprite: &entities.Sprite{
							Img: skeletonImg,
							X:   obj.X,
							Y:   obj.Y,
						},
						FollowsPlayer: follows,
						CombatComp:    components.NewEnemieCombat(3, 1, 60), // Cooldown de 1s (60 ticks)
					}
					g.enemies = append(g.enemies, newEnemy)

				case "training_dummy":
					dummyImg, _, err := ebitenutil.NewImageFromFile("./assets/images/dummy.png")
					if err != nil {
						log.Fatal(err)
					}
					newDummy := entities.NewTrainingDummy(obj.X, obj.Y, dummyImg)
					g.dummies = append(g.dummies, newDummy)

				case "potion_spawn":
					potionImg, _, err := ebitenutil.NewImageFromFile("./assets/images/health.png")
					if err != nil {
						log.Fatal(err)
					}
					amount, found := tilemap.GetIntProperty("amount", obj.Properties)
					if !found {
						amount = 2
					}
					newPotion := &entities.Potion{
						Sprite: &entities.Sprite{
							Img: potionImg,
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
