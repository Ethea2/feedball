package entities

import (
	"image"
	"math"

	"github.com/Ethea2/feedball/animation"
	"github.com/Ethea2/feedball/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState uint8

type PlayerFacing uint8

type PlayerSubstate uint8

// Enum for player states
const (
	Idle PlayerState = iota
	Running
	Down
	Jumping
)

const (
	SubIdle PlayerSubstate = iota
	FreeFalling
	HoldingTheBall
)

// Enum for where player is facing
const (
	Left PlayerFacing = iota
	Right
)

type AnimationKey struct {
	State  PlayerState
	Facing PlayerFacing
}

type PlayerInput struct {
	Left, Right, Jump, Down, Shoot ebiten.Key
}

type Player struct {
	*Sprite
	State               PlayerState
	Facing              PlayerFacing
	SubState            PlayerSubstate
	JumpTimer           float64
	Animations          map[AnimationKey]*animation.Animation
	movementRestricted  bool
	IsAffectedByGravity bool
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

func (p *Player) isGrounded(wallsAndFloors []image.Rectangle) bool {
	feetRect := image.Rect(
		int(math.Round(p.X)),
		int(math.Round(p.Y))+p.Bounds().Dy(), // bottom edge using Bounds
		int(math.Round(p.X))+p.Bounds().Dx(),
		int(math.Round(p.Y))+p.Bounds().Dy()+1, // one pixel below
	)
	for _, collider := range wallsAndFloors {
		if collider.Overlaps(feetRect) {
			return true
		}
	}
	return false
}

func (p *Player) Update(wallsAndFloors []image.Rectangle, playerInput PlayerInput) {
	p.Dx = 0.0

	if ebiten.IsKeyPressed(playerInput.Left) && !p.movementRestricted {
		p.Dx = -constants.PlayerSpeed
		p.Facing = Left
		if p.State != Jumping {
			p.State = Running
		}
	} else if ebiten.IsKeyPressed(playerInput.Right) && !p.movementRestricted {
		p.Dx = constants.PlayerSpeed
		p.Facing = Right
		if p.State != Jumping {
			p.State = Running
		}
	} else if p.State != Jumping && !p.movementRestricted {
		p.State = Idle
	}

	if ebiten.IsKeyPressed(playerInput.Jump) && p.SubState != FreeFalling {
		p.Jump()
	}

	if p.IsAffectedByGravity {
		p.Dy += constants.Gravity
	}

	p.X += p.Dx
	p.Sprite.CheckCollisionHorizontal(p.Sprite, wallsAndFloors)

	p.Y += p.Dy
	p.Sprite.CheckCollisionVertical(p.Sprite, wallsAndFloors)

	if p.isGrounded(wallsAndFloors) {
		if p.State == Jumping || p.SubState == FreeFalling {
			p.State = Idle
			p.movementRestricted = false
			p.IsAffectedByGravity = true
			if p.SubState == HoldingTheBall || p.SubState == FreeFalling {
				p.SubState = SubIdle
			}
		}
	}

	//Handle freefalling when player not holding ball and is jumping. Should not be accessed if player is already free falling
	if p.SubState != HoldingTheBall && p.State == Jumping && p.SubState != FreeFalling {
		if ebiten.IsKeyPressed(playerInput.Down) {
			p.SubState = FreeFalling
			p.movementRestricted = true
			p.IsAffectedByGravity = false
			p.Dy += constants.PlayerSpeed * 2
		}
	}
	p.UpdateAnimation()
}
