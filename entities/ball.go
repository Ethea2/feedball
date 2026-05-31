package entities

import (
	"image"

	"github.com/Ethea2/feedball/shared"
)

type BallState uint8

const (
	Initial BallState = iota
	PlayerHeld
	BallInPlay
)

type Ball struct {
	*Sprite
	State     BallState
	Speed     float64
	Direction shared.Direction
	ThrownBy  *Player
}

func (b *Ball) IsActive() bool {
	return b.State == BallInPlay
}

func (b *Ball) ReflectX() {
	switch b.Direction {
	case shared.East:
		b.Direction = shared.West
	case shared.West:
		b.Direction = shared.East
	case shared.NorthEast:
		b.Direction = shared.NorthWest
	case shared.NorthWest:
		b.Direction = shared.NorthEast
	case shared.SouthEast:
		b.Direction = shared.SouthWest
	case shared.SouthWest:
		b.Direction = shared.SouthEast
	}
}

func (b *Ball) ReflectY() {
	switch b.Direction {
	case shared.NorthEast:
		b.Direction = shared.SouthEast
	case shared.SouthEast:
		b.Direction = shared.NorthEast
	case shared.NorthWest:
		b.Direction = shared.SouthWest
	case shared.SouthWest:
		b.Direction = shared.NorthWest
	}
}

func (b *Ball) ResolveCollision(collider image.Rectangle) bool {
	if !b.Bounds().Overlaps(collider) {
		return false
	}

	overlapX := min(b.Bounds().Max.X, collider.Max.X) - max(b.Bounds().Min.X, collider.Min.X)
	overlapY := min(b.Bounds().Max.Y, collider.Max.Y) - max(b.Bounds().Min.Y, collider.Min.Y)

	if overlapX < overlapY {
		b.ReflectX()
		switch b.Direction {
		case shared.East, shared.NorthEast, shared.SouthEast:
			b.X += float64(overlapX)
		case shared.West, shared.NorthWest, shared.SouthWest:
			b.X -= float64(overlapX)
		}
	} else {
		b.ReflectY()
		switch b.Direction {
		case shared.SouthEast, shared.SouthWest:
			b.Y += float64(overlapY)
		case shared.NorthEast, shared.NorthWest:
			b.Y -= float64(overlapY)
		}
	}

	//make the ball catchable
	b.ThrownBy = nil
	return true
}

func (b *Ball) Update(wallsAndFloors []image.Rectangle) {
	if b.State != BallInPlay {
		return
	}
	switch b.Direction {
	case shared.East:
		b.X += b.Speed
	case shared.West:
		b.X -= b.Speed
	case shared.NorthEast:
		b.X += b.Speed
		b.Y -= b.Speed
	case shared.NorthWest:
		b.X -= b.Speed
		b.Y -= b.Speed
	case shared.SouthEast:
		b.X += b.Speed
		b.Y += b.Speed
	case shared.SouthWest:
		b.X -= b.Speed
		b.Y += b.Speed
	}
	for _, collider := range wallsAndFloors {
		if b.ResolveCollision(collider) {
			break
		}
	}
}
