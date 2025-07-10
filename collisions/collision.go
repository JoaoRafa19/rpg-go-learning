package collisions

import "image"

const CellSize = 64

type Grid struct {
	cols, rows int
	cells      [][]_Cells
}

type _Cells struct {
	colliders []*image.Rectangle
}

func NewGrid(width, height int) *Grid {
	cols := (width + CellSize - 1) / CellSize
	rows := (height + CellSize - 1) / CellSize

	cells := make([][]_Cells, rows)

	for i := range cells {
		cells[i] = make([]_Cells, cols)
		for j := range cells[i] {
			cells[i][j].colliders = make([]*image.Rectangle, 0)
		}
	}

	return &Grid{
		cols:  cols,
		rows:  rows,
		cells: cells,
	}
}

func (g *Grid) Insert(collider *image.Rectangle) {
	minX := collider.Min.X / CellSize
	maxX := (collider.Max.X - 1) / CellSize
	minY := collider.Min.Y / CellSize
	maxY := (collider.Max.Y - 1) / CellSize

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if x >= 0 && x < g.cols && y >= 0 && y < g.rows {
				g.cells[y][x].colliders = append(g.cells[y][x].colliders, collider)
			}
		}
	}
}

// GetNearbyColliders retorna todos os colisores únicos que estão próximos a uma área (bounds).
func (g *Grid) GetNearbyColliders(bounds image.Rectangle) []*image.Rectangle {
	nearby := make(map[*image.Rectangle]struct{}) // Usamos um map para evitar duplicatas

	minX := bounds.Min.X / CellSize
	maxX := (bounds.Max.X - 1) / CellSize
	minY := bounds.Min.Y / CellSize
	maxY := (bounds.Max.Y - 1) / CellSize

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if x >= 0 && x < g.cols && y >= 0 && y < g.rows {
				for _, collider := range g.cells[y][x].colliders {
					nearby[collider] = struct{}{}
				}
			}
		}
	}

	// Converte o map de volta para um slice
	result := make([]*image.Rectangle, 0, len(nearby))
	for collider := range nearby {
		result = append(result, collider)
	}
	return result
}
