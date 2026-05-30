package entities

import (
	"fmt"
	"image"

	"github.com/Ethea2/feedball/animation"
	"github.com/Ethea2/feedball/shared"
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerState uint8
type PlayerFacing uint8
type VisualState uint8

// Enum for player states
const (
	Idle PlayerState = iota
	Running
	Down
	Jumping
)

// Enum for visual states (animation only)
const (
	VisualIdle VisualState = iota
	VisualRunning
	VisualJumping
	VisualHoldingBallIdle
	VisualHoldingBallRunning
	VisualHoldingBallJumping
	VisualFreeFalling
	VisualThrowingWindup
	VisualThrowingNorthEast
	VisualThrowingEast
	VisualThrowingSouthEast
	VisualThrowingNorthWest
	VisualThrowingWest
	VisualThrowingSouthWest
)

// Enum for where player is facing
const (
	Left PlayerFacing = iota
	Right
)

type AnimationKey struct {
	State  VisualState
	Facing PlayerFacing
}

type PlayerInput struct {
	Left, Right, Jump, Down, Shoot ebiten.Key
}

type Player struct {
	*Sprite
	State               PlayerState
	Facing              PlayerFacing
	JumpTimer           float64
	Animations          map[AnimationKey]*animation.Animation
	ThrowAnimation      *animation.ThrowAnimation
	ThrowDirection      shared.Direction
	movementRestricted  bool
	IsAffectedByGravity bool
	FreeFalling         bool
	HoldingBall         bool
	Throwing            bool
}

func (p *Player) CurrentVisualState() VisualState {
	switch {
	case p.Throwing:
		return VisualThrowingWindup
	case p.FreeFalling:
		return VisualFreeFalling
	default:
		switch p.State {
		case Running:
			if p.HoldingBall {
				return VisualHoldingBallRunning
			}
			return VisualRunning
		case Jumping:
			if p.HoldingBall {
				return VisualHoldingBallJumping
			}
			return VisualJumping
		default:
			if p.HoldingBall {
				return VisualHoldingBallIdle
			}
			return VisualIdle
		}
	}
}

func (p *Player) ActiveThrowVisualState() VisualState {
	switch p.ThrowDirection {
	case shared.NorthEast:
		return VisualThrowingNorthEast
	case shared.SouthEast:
		return VisualThrowingSouthEast
	case shared.NorthWest:
		return VisualThrowingNorthWest
	case shared.SouthWest:
		return VisualThrowingSouthWest
	case shared.West:
		return VisualThrowingWest
	default:
		return VisualThrowingEast
	}
}

func (p *Player) Jump() {
	if p.State == Jumping {
		return
	}
	p.State = Jumping
	p.JumpTimer = 1.0
	p.Dy = -15.0
}

func (p *Player) StartThrow(windupFirst, windupLast int) {
	p.Throwing = true
	p.movementRestricted = true
	p.IsAffectedByGravity = false
	p.Dy = 0
	p.Dx = 0
	p.ThrowAnimation = animation.NewThrowAnimation(
		windupFirst, windupLast,
		0, 0, // release frames set later via SetRelease
		shared.ThrowFrameSpeed,
	)
}

func (p *Player) SampleThrowDirection(playerInput PlayerInput) shared.Direction {
	up := ebiten.IsKeyPressed(playerInput.Jump)
	down := ebiten.IsKeyPressed(playerInput.Down)

	if p.Facing == Right {
		switch {
		case up:
			return shared.NorthEast
		case down:
			return shared.SouthEast
		default:
			return shared.East
		}
	} else {
		switch {
		case up:
			return shared.NorthWest
		case down:
			return shared.SouthWest
		default:
			return shared.West
		}
	}
}

func (p *Player) FinishThrow() shared.Direction {
	dir := p.ThrowDirection
	p.Throwing = false
	p.HoldingBall = false
	p.movementRestricted = false
	p.IsAffectedByGravity = true
	p.ThrowAnimation = nil
	return dir
}

func (p *Player) ActiveAnimation() *animation.Animation {
	if p.Throwing {
		return nil
	}
	key := AnimationKey{State: p.CurrentVisualState(), Facing: p.Facing}
	if anim, ok := p.Animations[key]; ok {
		return anim
	}
	fallback := AnimationKey{State: VisualIdle, Facing: p.Facing}
	return p.Animations[fallback]
}

func (p *Player) UpdateAnimation(playerInput PlayerInput, releaseFrames map[shared.Direction][2]int) {
	if p.Throwing && p.ThrowAnimation != nil {
		if p.ThrowAnimation.IsOnLastWindupFrame() {
			p.ThrowDirection = p.SampleThrowDirection(playerInput)
			frames := releaseFrames[p.ThrowDirection]
			p.ThrowAnimation.SetRelease(frames[0], frames[1])
		}
		p.ThrowAnimation.Update()
		return
	}
	anim := p.ActiveAnimation()
	if anim != nil {
		anim.Update()
	}
}

func (p *Player) isGrounded(wallsAndFloors []image.Rectangle) bool {
	bounds := p.Bounds()
	feetRect := image.Rect(
		bounds.Min.X,
		bounds.Max.Y,
		bounds.Max.X,
		bounds.Max.Y+1,
	)
	for _, collider := range wallsAndFloors {
		if collider.Overlaps(feetRect) {
			return true
		}
	}
	return false
}

func (p *Player) Update(wallsAndFloors []image.Rectangle, playerInput PlayerInput, releaseFrames map[shared.Direction][2]int) {
	fmt.Println(p.CurrentVisualState())
	p.Dx = 0.0
	if ebiten.IsKeyPressed(playerInput.Left) && !p.movementRestricted {
		p.Dx = -shared.PlayerSpeed
		p.Facing = Left
		if p.State != Jumping {
			p.State = Running
		}
	} else if ebiten.IsKeyPressed(playerInput.Right) && !p.movementRestricted {
		p.Dx = shared.PlayerSpeed
		p.Facing = Right
		if p.State != Jumping {
			p.State = Running
		}
	} else if p.State != Jumping && !p.movementRestricted {
		p.State = Idle
	}

	if ebiten.IsKeyPressed(playerInput.Jump) && !p.FreeFalling {
		p.Jump()
	}

	if p.IsAffectedByGravity {
		p.Dy += shared.Gravity
	}

	p.X += p.Dx
	p.Sprite.CheckCollisionHorizontal(p.Sprite, wallsAndFloors)
	p.Y += p.Dy
	p.Sprite.CheckCollisionVertical(p.Sprite, wallsAndFloors)

	if p.State == Jumping || p.FreeFalling {
		if p.isGrounded(wallsAndFloors) {
			p.State = Idle
			p.movementRestricted = false
			p.IsAffectedByGravity = true
			if p.FreeFalling {
				p.FreeFalling = false
			}
		}
	}

	if !p.HoldingBall && p.State == Jumping && !p.FreeFalling {
		if ebiten.IsKeyPressed(playerInput.Down) {
			p.FreeFalling = true
			p.movementRestricted = true
			p.IsAffectedByGravity = false
			p.Dy += shared.PlayerSpeed * 2
		}
	}

	p.UpdateAnimation(playerInput, releaseFrames)
}
