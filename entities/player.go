package entities

import (
	"math"
	"rpg-go/animations"
	"rpg-go/camera"
	"rpg-go/components"
	"rpg-go/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState uint8

const (
	Down PlayerState = iota
	Up
	Left
	Right
)

type Player struct {
	*Sprite
	Animations map[PlayerState]*animations.Animation
	CombatComp *components.BasicCombat
}

func (p *Player) ActiveAnimation(dx, dy int) *animations.Animation {
	if dx > 0 {
		return p.Animations[Right]
	}
	if dx < 0 {
		return p.Animations[Left]
	}
	if dy > 0 {
		return p.Animations[Down]
	}
	if dy < 0 {
		return p.Animations[Up]
	}

	return nil

}

func (p *Player) GetY() float64 {
	return p.Y
}

func (p *Player) Draw(screen *ebiten.Image, cam *camera.Camera, sheet *spritesheet.SpriteSheet) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	opts.GeoM.Translate(cam.X, cam.Y)

	activeAnimation := p.ActiveAnimation(int(p.Dx), int(p.Dy))
	playerFrame := 0
	if activeAnimation != nil {
		// Pega o índice do frame atual da animação
		playerFrame = activeAnimation.Frame()
	} else {
		// Se estiver parado, pegue o primeiro frame da animação "Down" como padrão
		// Você pode ajustar isso para a direção que o jogador estava olhando por último
		playerFrame = p.Animations[Down].GetFirstFrame()
	}

	// Calcula o retângulo (área) do frame no spritesheet
	frameRect := sheet.Rect(playerFrame)

	opts.GeoM.Reset()
	opts.GeoM.Translate(p.X, p.Y)
	opts.GeoM.Translate(cam.X, cam.Y)

	// Desenha APENAS a sub-imagem correspondente ao frame
	screen.DrawImage(p.Img.SubImage(frameRect).(*ebiten.Image), opts)
}

func (p *Player) Move() {
	// React to keyPresses
	p.Dx = 0.0
	p.Dy = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Dy = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.Dx += 2
	}
	// Normalize movement
	// Magnitude
	magnitude := math.Sqrt(p.Dx*p.Dx + p.Dy*p.Dy)

	if magnitude > 0.0 {
		speed := 2.0
		p.Dx = (p.Dx / magnitude) * speed
		p.Dy = (p.Dy / magnitude) * speed
	}

}
