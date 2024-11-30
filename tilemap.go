package main

import (
	"encoding/json"
	"image"
	"os"
	"path"
	"strings"
)

// data we want for one layer in our list of layers
type TilemapLayerJSON struct {
	Data    []int      `json:"data"`
	Width   int        `json:"width"`
	Height  int        `json:"height"`
	Name    string     `json:"name"`
	Class   string     `json:"class"`
	Objects []Collisor `json:"objects"`
}

type Collisor struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Rect   image.Rectangle
}

// all layers in a tilemap
type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
	// raw data for each tileset (path, gid)
	Tilesets []map[string]any `json:"tilesets"`
}

func (t *TilemapJSON) GenCollisors() []Collisor {
	colliders := make([]Collisor, 0)

	for _, layer := range t.Layers {
		if strings.Compare(layer.Class, "obj") == 0 {
			for _, object := range layer.Objects {
				newCollider := Collisor{
					X:      object.X,
					Y:      object.Y,
					Width:  object.Width,
					Height: object.Height,
				} 
				newCollider.Rect = image.Rect(
					int(newCollider.X),
					int(newCollider.Y),
					int(newCollider.X)+int(newCollider.Width),
					int(newCollider.Y)+int(newCollider.Height),
				)
				colliders = append(colliders, newCollider)
			}
		}
	}

	return colliders
}

// temp function to generate all of our tilesets and return a slice of them
func (t *TilemapJSON) GenTilesets() ([]Tileset, error) {
	tilesets := make([]Tileset, 0)

	for _, tilesetData := range t.Tilesets {
		// convert map relative path to project relative path
		tilesetPath := path.Join("assets/maps/", tilesetData["source"].(string))
		tileset, err := NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}

		tilesets = append(tilesets, tileset)
	}

	return tilesets, nil
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
