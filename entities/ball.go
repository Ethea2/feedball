package entities

type BallState uint8

const (
	Initial BallState = iota
	PlayerHeld
	PlayerThrown
)

type Ball struct {
	*Sprite
}
