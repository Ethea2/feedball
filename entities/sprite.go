package entities

import (
	"image"
	"math"

	"github.com/Ethea2/feedball/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img              *ebiten.Image
	X, Y, Dx, Dy     float64
	XOffset, YOffset int
}

func (s *Sprite) CheckCollisionHorizontal(sprite *Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(sprite.Bounds()) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(
					collider.Min.X,
				) - float64(
					constants.TileSize*constants.PlayerBallScale,
				)
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
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
					constants.TileSize*constants.PlayerBallScale,
				)
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
			sprite.Dy = 0
		}
	}
}

func (s *Sprite) Bounds() image.Rectangle {
	return image.Rect(
		int(math.Round(s.X)),
		int(math.Round(s.Y)),
		int((math.Round(s.X))+constants.TileSize*2)+s.XOffset,
		int((math.Round(s.Y))+constants.TileSize*2)+s.YOffset,
	)
}
