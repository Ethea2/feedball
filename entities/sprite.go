package entities

import (
	"image"
	"math"

	"github.com/Ethea2/feedball/shared"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img                                               *ebiten.Image
	X, Y, Dx, Dy                                      float64
	LeftXOffset, UpYOffset, RightXOffset, DownYOffset int
}

func (s *Sprite) CheckCollisionHorizontal(sprite *Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(sprite.Bounds()) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(
					collider.Min.X,
				) - float64(
					shared.TileSize*2,
				) - float64(
					sprite.RightXOffset,
				)
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X) - float64(sprite.LeftXOffset)
			}
			sprite.Dx = 0
		}
	}
}

func (s *Sprite) CheckCollisionVertical(sprite *Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(sprite.Bounds()) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(
					collider.Min.Y,
				) - float64(
					shared.TileSize*2,
				) - float64(
					sprite.DownYOffset,
				)
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y) - float64(sprite.UpYOffset)
			}
			sprite.Dy = 0
		}
	}
}

func (s *Sprite) Bounds() image.Rectangle {
	return image.Rect(
		int(math.Round(s.X))+s.LeftXOffset,
		int(math.Round(s.Y))+s.UpYOffset,
		int((math.Round(s.X))+shared.TileSize*2)+s.RightXOffset,
		int((math.Round(s.Y))+shared.TileSize*2)+s.DownYOffset,
	)
}
