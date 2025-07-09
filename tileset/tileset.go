package tileset

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"rpg-go/constants"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TilesetJSON espelha a estrutura de um arquivo .tsx exportado como .json
type TilesetJSON struct {
	Image       string     `json:"image"`       // Usado por tilesets baseados em uma única imagem (spritesheet)
	Columns     int        `json:"columns"`     // Número de colunas no spritesheet. ESSENCIAL!
	TileWidth   int        `json:"tilewidth"`
	TileHeight  int        `json:"tileheight"`
	Tiles       []TileJSON `json:"tiles"`       // Usado por tilesets baseados em uma coleção de imagens
}

// TileJSON representa um único tile dentro de uma coleção de imagens.
type TileJSON struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
}

// Tileset é a nossa estrutura unificada. Ela pode representar tanto um
// tileset de imagem única (spritesheet) quanto uma coleção de imagens.
type Tileset struct {
	FirstGid int
	
	// Para tilesets de imagem única (spritesheet)
	spritesheet *ebiten.Image
	columns     int

	// Para tilesets de coleção de imagens
	individualTiles map[int]*ebiten.Image
}

// NewTileset é a nossa factory. Ela lê um arquivo de tileset do Tiled,
// determina seu tipo, e retorna uma struct Tileset pronta para uso.
func NewTileset(path string, firstGid int) (*Tileset, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o arquivo do tileset %s: %w", path, err)
	}

	var data TilesetJSON
	if err := json.Unmarshal(contents, &data); err != nil {
		return nil, fmt.Errorf("falha ao decodificar o JSON do tileset %s: %w", path, err)
	}

	tileset := &Tileset{
		FirstGid: firstGid,
	}

	// baseDir é o diretório onde o arquivo .tsx/.json está, para resolver caminhos relativos.
	baseDir := filepath.Dir(path)

	// DETERMINAÇÃO DE TIPO:
	// Se a propriedade "image" existe, é um tileset de imagem única (spritesheet).
	if data.Image != "" {
		// É um tileset de imagem única (Uniform)
		imgPath := filepath.Join(baseDir, data.Image)
		img, _, err := ebitenutil.NewImageFromFile(imgPath)
		if err != nil {
			return nil, fmt.Errorf("falha ao carregar a imagem do tileset %s: %w", imgPath, err)
		}
		tileset.spritesheet = img
		tileset.columns = data.Columns // Lendo o número de colunas do arquivo!
		
		if tileset.columns == 0 {
			log.Printf("Aviso: Tileset '%s' não tem 'columns' definido. Assumindo que a largura é a da imagem.", path)
			tileset.columns = img.Bounds().Dx() / constants.Tilesize
		}
		
	} else {
		// É um tileset de coleção de imagens (Dynamic)
		tileset.individualTiles = make(map[int]*ebiten.Image)
		for _, tileData := range data.Tiles {
			imgPath := filepath.Join(baseDir, tileData.Image)
			img, _, err := ebitenutil.NewImageFromFile(imgPath)
			if err != nil {
				return nil, fmt.Errorf("falha ao carregar a imagem do tile individual %s: %w", imgPath, err)
			}
			tileset.individualTiles[tileData.ID] = img
		}
	}

	return tileset, nil
}

func (t *Tileset) Img(id int) *ebiten.Image {
	localID := id - t.FirstGid

	if t.individualTiles != nil {
		return t.individualTiles[localID]
	}

	if t.spritesheet != nil {
		srcX := (localID % t.columns) * constants.Tilesize
		srcY := (localID / t.columns) * constants.Tilesize

		rect := image.Rect(srcX, srcY, srcX+constants.Tilesize, srcY+constants.Tilesize)
		return t.spritesheet.SubImage(rect).(*ebiten.Image)
	}

	return nil 
}