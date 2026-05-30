package entities

import "github.com/Ethea2/feedball/shared"

type BallState uint8

const (
	Initial BallState = iota
	PlayerHeld
	PlayerThrown
)

type Ball struct {
	*Sprite
	State     BallState
	Speed     float64
	Direction shared.Direction
}

func (b *Ball) IsActive() bool {
	return b.State == PlayerThrown
}

func (b *Ball) Update() {
	if b.State != PlayerThrown {
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
}
