package tilemap

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"rpg-go/tileset"
)

// data we want for one layer in our list of layers
type TilemapLayerJSON struct {
	Data       []int         `json:"data"`
	Width      int           `json:"width"`
	Height     int           `json:"height"`
	Name       string        `json:"name"`
	Type       string        `json:"type"` // "tilelayer" ou "objectgroup"
	Objects    []TiledObject `json:"objects"`
	Collisions []TiledObject `json:"collisions"`
}

type TiledProperty struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type TiledObject struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	X          float64         `json:"x"`
	Y          float64         `json:"y"`
	Width      float64         `json:"width"`
	Height     float64         `json:"height"`
	GID        int             `json:"gid,omitempty"`
	Properties []TiledProperty `json:"properties"`
}

// all layers in a tilemap
type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
	// raw data for each tileset (path, gid)
	Tilesets []map[string]any `json:"tilesets"`
}

// temp function to generate all of our tilesets and return a slice of them
func (t *TilemapJSON) GenTilesets() ([]*tileset.Tileset, error) {
	tilesets := make([]*tileset.Tileset, 0)

	for _, tilesetData := range t.Tilesets {
		// convert map relative path to project relative path
		tilesetPath := path.Join("assets/maps/", tilesetData["source"].(string))
		tileset, err := tileset.NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}

		tilesets = append(tilesets, tileset)
	}

	return tilesets, nil
}

func GetIntProperty(name string, properties []TiledProperty) (int, bool) {
	for _, prop := range properties {
		if prop.Name == name {
			// O Tiled exporta números como float64, então precisamos converter.
			if value, ok := prop.Value.(float64); ok {
				return int(value), true
			}
			log.Printf("Aviso: Propriedade '%s' encontrada, mas não é um número (float64).", name)
			return 0, false
		}
	}
	return 0, false // Retorna o valor padrão 0 se a propriedade não for encontrada.
}

// opens the file, parses it, and returns the json object + potential error
func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
