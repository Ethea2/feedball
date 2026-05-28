package entities

type BallState uint8

const (
	Initial BallState = iota
	PlayerHeld
	PlayerThrown
)

type Ball struct {
	*Sprite
	Active bool //ball in play or not? I'M NOT AI MOTHAFUCKA
}
