package tileset

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"rpg-go/constants"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
	Gid  int
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid
	// get the position of the image where the tile id is
	srcX := (id) % 22
	srcY := (id) / 22

	// convert to pixel position
	srcX *= constants.Tilesize
	srcY *= constants.Tilesize

	return u.img.SubImage(image.Rect(srcX, srcY, srcX+constants.Tilesize, srcY+constants.Tilesize)).(*ebiten.Image)

}

type TilesetType int

const (
	Uniform TilesetType = iota
	Dynamic
)

type TilesetJSONProperty struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value TilesetType `json:"value"`
}

type TilesetJSON struct {
	Properties []*TilesetJSONProperty `json:"properties"`
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynamicTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynTileset) Img(id int) *ebiten.Image {
	id -= d.gid

	return d.imgs[id]
}

func normalizePath(path string) string {
	op := path
	op = filepath.Clean(op)
	op = strings.ReplaceAll(op, "\\", "/")
	op = strings.TrimPrefix(op, "../")
	op = strings.TrimPrefix(op, "../")
	op = filepath.Join("./assets/", op)
	fmt.Println(op)
	return op
}

func NewTileset(path string, git int) (Tileset, error) {
	contents, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Load file err", err)
		return nil, err
	}

	var tilesetJSON TilesetJSON

	err = json.Unmarshal(contents, &tilesetJSON)
	if err != nil {
		fmt.Println("Unmarshal file err", err)
		return nil, err
	}
	var t TilesetType
	if strings.Contains(path, "buildings") {
		t = Dynamic
	} else {
		t = Uniform
	}

	switch t {
	case Uniform:
		uniTileJson := &UniformTilesetJSON{}
		err = json.Unmarshal(contents, uniTileJson)
		if err != nil {
			fmt.Println("Unmarshal content err", err)
			return nil, err
		}

		uniformTileset := &UniformTileset{}

		img, _, err := ebitenutil.NewImageFromFile(normalizePath(uniTileJson.Path))
		if err != nil {
			fmt.Println("EbitenImage err", err)
			return nil, err
		}
		uniformTileset.img = img
		uniformTileset.gid = git
		return uniformTileset, nil

	case Dynamic:
		dynTilesJson := &DynamicTilesetJSON{}
		err = json.Unmarshal(contents, dynTilesJson)
		if err != nil {
			return nil, err
		}
		var dynTilset DynTileset
		dynTilset.gid = git
		dynTilset.imgs = make([]*ebiten.Image, 0)
		for _, tileJson := range dynTilesJson.Tiles {
			imgPath := filepath.Join("maps", "tilesets", tileJson.Path)
			imgPath = normalizePath(imgPath)
			imgPath = strings.ReplaceAll(imgPath, "\\", "/")
			img, _, err := ebitenutil.NewImageFromFile(imgPath)
			if err != nil {
				return nil, err
			}
			dynTilset.imgs = append(dynTilset.imgs, img)
		}
		return &dynTilset, nil
	}

	// return uniform
	return nil, fmt.Errorf("Tileset not found in %s", path)
}
