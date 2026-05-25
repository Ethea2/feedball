package entities

type PlayerState uint8

const (
	Down PlayerState = iota
	Up
	Left
	Right
	Jumping
)

type Player struct {
	*Sprite
	State     PlayerState
	JumpTimer float64
}

func (p *Player) Jump() {
	p.State = Jumping
	p.JumpTimer = 1.0
}
