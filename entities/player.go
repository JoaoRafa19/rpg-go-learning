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
	AttackRight
	AttackLeft
)

type Player struct {
	*Sprite
	Animations map[PlayerState]*animations.Animation
	CombatComp *components.BasicCombat

	Facing      PlayerState
	isAttacking bool // Atacando agora?
	AttackTick  int  // Duração do ataque
}

func NewPlayer(img *ebiten.Image) *Player {
	return &Player{
		Animations: map[PlayerState]*animations.Animation{
			Up:          animations.NewAnimation(5, 13, 4, 20.0),
			Down:        animations.NewAnimation(4, 12, 4, 20),
			Left:        animations.NewAnimation(6, 14, 4, 20),
			Right:       animations.NewAnimation(7, 15, 4, 20),
			AttackRight: animations.NewAnimation(19, 19, 0, 20),
			AttackLeft:  animations.NewAnimation(18, 18, 0, 20),
		},
		Facing: Down,

		CombatComp: components.NewBasicCombat(10, 1), // Aumentei a vida para 10
		Sprite: &Sprite{
			Img: img,
		},
	}
}

func (p *Player) ActiveAnimation() *animations.Animation {

	if p.isAttacking {

		return p.Animations[AttackRight]
	}

	// Lógica de movimento (como antes)
	if p.Dx > 0 {
		return p.Animations[Right]
	}
	if p.Dx < 0 {
		return p.Animations[Left]
	}
	if p.Dy > 0 {
		return p.Animations[Down]
	}
	if p.Dy < 0 {
		return p.Animations[Up]
	}

	return nil // Parado

}

func (p *Player) GetY() float64 {
	return p.Y
}

func (p *Player) Draw(screen *ebiten.Image, cam *camera.Camera, sheet *spritesheet.SpriteSheet) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	opts.GeoM.Translate(cam.X, cam.Y)

	activeAnimation := p.ActiveAnimation()
	playerFrame := 0
	if activeAnimation != nil {
		// Pega o índice do frame atual da animação
		playerFrame = activeAnimation.Frame()
	} else {
		// Se estiver parado, pegue o primeiro frame da animação "Down" como padrão
		// Você pode ajustar isso para a direção que o jogador olhou por último
		playerFrame = p.Animations[Down].GetFirstFrame()
	}

	// Calcula o retângulo (área) do frame no spritesheet
	frameRect := sheet.Rect(playerFrame)

	opts.GeoM.Reset()
	opts.GeoM.Translate(p.X, p.Y)
	opts.GeoM.Translate(cam.X, cam.Y)

	// Desenha APENAS a imagem correspondente ao frame
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

const playerAttackDuration = 20

func (p *Player) IsAttacking() bool {
	return p.isAttacking
}

func (p *Player) Attack() {
	if !p.isAttacking {
		p.isAttacking = true
		p.AttackTick = 0
		// Reinicia a animação de ataque correspondente
		// p.Animations[p.facingAttackState()].Reset() // Veremos isso depois
	}
}

func (p *Player) SetFacing(dir PlayerState) {
	p.Facing = dir
}

func (p *Player) UpdateAttack() {
	if p.isAttacking {
		p.AttackTick++
		if p.AttackTick >= playerAttackDuration {
			p.isAttacking = false // O tempo do ataque acabou, voltamos ao estado normal.
		}
	}
}

func (p *Player) UpdateAttackTick() {
	if p.isAttacking {
		p.AttackTick++
		if p.AttackTick >= playerAttackDuration {
			p.isAttacking = false
			p.AttackTick = 0
		}
	}
}
