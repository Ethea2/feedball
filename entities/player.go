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

func (p *Player) checkGrounded(wallsAndFloors []image.Rectangle) bool {
	// Check one pixel below the player's feet
	feetRect := image.Rect(
		int(math.Round(p.X)),
		int(math.Round(p.Y))+constants.TileSize, // bottom edge
		int(math.Round(p.X))+constants.TileSize,
		int(math.Round(p.Y))+constants.TileSize+1, // one pixel below
	)
	for _, collider := range wallsAndFloors {
		if collider.Overlaps(feetRect) {
			return true
		}
	}
	return false
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
		p.State = Idle
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Jump()
	}

	p.Dy += constants.Gravity

	p.X += p.Dx
	p.Sprite.CheckCollisionHorizontal(p.Sprite, wallsAndFloors)

	p.Y += p.Dy

	wasGrounded := p.checkGrounded(wallsAndFloors)
	p.Sprite.CheckCollisionVertical(p.Sprite, wallsAndFloors)

	if p.State == Jumping && wasGrounded && p.Dy >= 0 {
		p.State = Idle
	}

	p.UpdateAnimation()
}
