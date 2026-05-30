package entities

import "github.com/Ethea2/feedball/shared"

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

func (b *Ball) Update() {
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
}
