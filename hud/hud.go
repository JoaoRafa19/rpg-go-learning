package hud

import (
	"fmt"
	"image/color"
	"rpg-go/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type HUD struct {
	playerCombat components.Combat
}

func NewHUD(playerCombat components.Combat) *HUD {
	return &HUD{
		playerCombat: playerCombat,
	}
}

func (h *HUD) Draw(screen *ebiten.Image) {
	barWidth := 50
	barHeight := 8

	px := float32(screen.Bounds().Dx()) - 55

	// Fundo da barra (vermelho escuro)
	vector.DrawFilledRect(screen, px, 5, float32(barWidth), float32(barHeight), color.RGBA{100, 0, 0, 255}, false)

	// Vida atual (verde)
	healthPercentage := float32(h.playerCombat.Health()) / float32(h.playerCombat.MaxHealth())
	currentHealthWidth := float32(barWidth) * healthPercentage
	vector.DrawFilledRect(screen, px, 5, currentHealthWidth, float32(barHeight), color.RGBA{0, 255, 0, 255}, false)

	// Texto da vida
	healthText := fmt.Sprintf("%d / %d", h.playerCombat.Health(), h.playerCombat.MaxHealth())
	ebitenutil.DebugPrintAt(screen, healthText, int(px), 5)
}
