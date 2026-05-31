package spritesheet

import "image"

type SpriteSheet struct {
	WidthInTiles  int
	HeightInTiles int
	Tilesize      int
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := (index % s.WidthInTiles) * s.Tilesize
	y := (index / s.WidthInTiles) * s.Tilesize

	return image.Rect(
		x, y, x+s.Tilesize, y+s.Tilesize,
	)
}

func (s *SpriteSheet) Rect64(index int) image.Rectangle {
	x := (index % s.WidthInTiles) * (s.Tilesize * 2)
	y := (index / s.WidthInTiles) * (s.Tilesize * 2)

	return image.Rect(
		x, y, x+(s.Tilesize*2), y+(s.Tilesize*2),
	)
}

func NewSpriteSheet(w, h, t int) *SpriteSheet {
	return &SpriteSheet{
		w, h, t,
	}
}
