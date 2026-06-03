package entities

import (
	"image"
	"math"

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

// In shared/directions.go or similar
func DirectionToVelocity(dir shared.Direction, speed float64) (dx, dy float64) {
	switch dir {
	case shared.East:
		return speed, 0
	case shared.West:
		return -speed, 0
	case shared.NorthEast:
		return speed, -speed
	case shared.NorthWest:
		return -speed, -speed
	case shared.SouthEast:
		return speed, speed
	case shared.SouthWest:
		return -speed, speed
	}
	return 0, 0
}

func (b *Ball) ResolveCollision(collider image.Rectangle) bool {
	if !b.Bounds().Overlaps(collider) {
		return false
	}

	overlapX := min(b.Bounds().Max.X, collider.Max.X) - max(b.Bounds().Min.X, collider.Min.X)
	overlapY := min(b.Bounds().Max.Y, collider.Max.Y) - max(b.Bounds().Min.Y, collider.Min.Y)

	dirBefore := b.Direction
	if overlapX < overlapY {
		b.ReflectX()
		switch dirBefore {
		case shared.East, shared.NorthEast, shared.SouthEast:
			b.X -= float64(overlapX) // was going right, push left
		case shared.West, shared.NorthWest, shared.SouthWest:
			b.X += float64(overlapX) // was going left, push right
		}
	} else {
		b.ReflectY()
		switch dirBefore {
		case shared.SouthEast, shared.SouthWest:
			b.Y -= float64(overlapY) // was going down, push up
		case shared.NorthEast, shared.NorthWest:
			b.Y += float64(overlapY) // was going up, push down
		}
	}

	//make the ball catchable
	b.ThrownBy = nil
	return true
}

// CLAUDE ONE SHOT THIS MOTHERFUCKING FEATURE. I KNEW THE SOLUTION AND WAS ASKING CLAUDE IF THE SWEPT AABB WOULD SOLVE MY PROBLEM AND BRO LAID OUT THE CODE/SOLUTION IN FRONT OF ME
// I MEAN TO BE FUCKING FAIR LIKE GODDAMN WHY IS IT SO GOOD T-T
// I FEEL LIKE SHIT. I'LL IMPLEMENT EVERYTHING ELSE ON MY OWN FROM HERE... FUCK YOU CLAUDE (FOR BEING TOO GOOD)
// SweptAABB returns time of impact [0,1] and collision normal.
// Returns 1.0, (0,0) if no collision this frame.
func SweptAABB(movingBounds image.Rectangle, dx, dy float64, staticBounds image.Rectangle) (tHit float64, normalX, normalY float64) {
	// Expand static bounds by moving bounds size (Minkowski sum)
	expandedMin := image.Point{
		X: staticBounds.Min.X - movingBounds.Dx(),
		Y: staticBounds.Min.Y - movingBounds.Dy(),
	}
	expandedMax := staticBounds.Max

	// Ray origin is the moving rect's Min corner
	originX := float64(movingBounds.Min.X)
	originY := float64(movingBounds.Min.Y)

	// Time to hit each face
	var tEntryX, tEntryY, tExitX, tExitY float64

	if dx != 0 {
		tEntryX = (float64(expandedMin.X) - originX) / dx
		tExitX = (float64(expandedMax.X) - originX) / dx
		if tEntryX > tExitX {
			tEntryX, tExitX = tExitX, tEntryX
		}
	} else {
		if originX < float64(expandedMin.X) || originX > float64(expandedMax.X) {
			return 1.0, 0, 0 // no collision
		}
		tEntryX, tExitX = math.Inf(-1), math.Inf(1)
	}

	if dy != 0 {
		tEntryY = (float64(expandedMin.Y) - originY) / dy
		tExitY = (float64(expandedMax.Y) - originY) / dy
		if tEntryY > tExitY {
			tEntryY, tExitY = tExitY, tEntryY
		}
	} else {
		if originY < float64(expandedMin.Y) || originY > float64(expandedMax.Y) {
			return 1.0, 0, 0
		}
		tEntryY, tExitY = math.Inf(-1), math.Inf(1)
	}

	tEntry := math.Max(tEntryX, tEntryY)
	tExit := math.Min(tExitX, tExitY)

	// No collision if ranges don't overlap or collision is behind/beyond movement
	if tEntry > tExit || tEntry < 0 || tEntry > 1.0 {
		return 1.0, 0, 0
	}

	// Determine normal from which axis was hit first
	if tEntryX > tEntryY {
		if dx < 0 {
			return tEntry, 1, 0
		}
		return tEntry, -1, 0
	}
	if dy < 0 {
		return tEntry, 0, 1
	}
	return tEntry, 0, -1
}

func (b *Ball) Update(wallsAndFloors []image.Rectangle) {
	if b.State != BallInPlay {
		return
	}

	dx, dy := DirectionToVelocity(b.Direction, b.Speed)

	earliestT := 1.0
	var hitNormalX, hitNormalY float64

	for _, collider := range wallsAndFloors {
		t, nx, ny := SweptAABB(b.Bounds(), dx, dy, collider)
		if t < earliestT {
			earliestT = t
			hitNormalX = nx
			hitNormalY = ny
		}
	}

	// Move to point of impact
	b.X += dx * earliestT
	b.Y += dy * earliestT

	// Reflect direction if we hit something
	if earliestT < 1.0 {
		b.reflectFromNormal(hitNormalX, hitNormalY)
		b.ThrownBy = nil
	}
}

func (b *Ball) reflectFromNormal(nx, ny float64) {
	switch {
	case nx != 0: // hit a vertical wall
		b.ReflectX()
	case ny != 0: // hit a horizontal wall
		b.ReflectY()
	}
}
