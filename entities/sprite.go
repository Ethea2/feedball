package entities

import (
	"image"
	"math"

	"github.com/Ethea2/feedball/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img          *ebiten.Image
	X, Y, Dx, Dy float64
}

func (s *Sprite) CheckCollisionHorizontal(sprite *Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(math.Round(sprite.X)),
				int(math.Round(sprite.Y)),
				int(math.Round(sprite.X))+constants.TileSize,
				int(math.Round(sprite.Y))+constants.TileSize,
			),
		) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - constants.TileSize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
			sprite.Dx = 0
		}
	}
}

func (s *Sprite) CheckCollisionVertical(sprite *Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(math.Round(sprite.X)),
				int(math.Round(sprite.Y)),
				int(math.Round(sprite.X))+constants.TileSize,
				int(math.Round(sprite.Y))+constants.TileSize,
			),
		) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - constants.TileSize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
			sprite.Dy = 0
		}
	}
}
