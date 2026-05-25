package entities

import (
	"image"

	"github.com/Ethea2/feedball/animation"
	"github.com/Ethea2/feedball/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState uint8

type PlayerFacing uint8

const (
	Idle PlayerState = iota
	Running
	Down
	Jumping
	HoldingTheBall
)

const (
	Left PlayerFacing = iota
	Right
)

type AnimationKey struct {
	State  PlayerState
	Facing PlayerFacing
}

type Player struct {
	*Sprite
	State      PlayerState
	Facing     PlayerFacing
	JumpTimer  float64
	Animations map[AnimationKey]*animation.Animation
}

func (p *Player) Jump() {
	if p.State == Jumping {
		return
	}
	p.State = Jumping
	p.JumpTimer = 1.0
	p.Dy = -15.0
}

func (p *Player) ActiveAnimation() *animation.Animation {
	key := AnimationKey{State: p.State, Facing: p.Facing}
	if anim, ok := p.Animations[key]; ok {
		return anim
	}

	fallback := AnimationKey{State: Idle, Facing: p.Facing}
	return p.Animations[fallback]
}

func (p *Player) UpdateAnimation() {
	anim := p.ActiveAnimation()
	if anim != nil {
		anim.Update()
	}
}

func (p *Player) Update(wallsAndFloors []image.Rectangle) {
	p.Dx = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Dx = -constants.PlayerSpeed
		p.Facing = Left
		if p.State != Jumping {
			p.State = Running
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Dx = constants.PlayerSpeed
		p.Facing = Right
		if p.State != Jumping {
			p.State = Running
		}
	} else if p.State != Jumping {
		p.State = Idle // reset to idle when no keys pressed
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Jump()
	}

	if p.State == Jumping {
		p.JumpTimer -= 0.02
		if p.JumpTimer <= 0 {
			p.State = Idle
			p.JumpTimer = 0
		}
	}

	p.Dy += constants.Gravity

	p.X += p.Dx
	p.Sprite.CheckCollisionHorizontal(p.Sprite, wallsAndFloors)

	p.Y += p.Dy
	p.Sprite.CheckCollisionVertical(p.Sprite, wallsAndFloors)

	p.UpdateAnimation()
}
