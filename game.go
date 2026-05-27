package main

import (
	"image"
	"image/color"
	"log"

	"github.com/Ethea2/feedball/animation"
	"github.com/Ethea2/feedball/constants"
	"github.com/Ethea2/feedball/entities"
	"github.com/Ethea2/feedball/spritesheet"
	"github.com/Ethea2/feedball/tiles"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	tilemapJSON       *tiles.TilemapJSON
	tilesets          []tiles.Tileset
	tilemapImg        *ebiten.Image
	wallsAndFloors    []image.Rectangle
}

func NewGame() *Game {
	tilemapImg, _, err := ebitenutil.NewImageFromFile("./assets/images/terrain.png")

	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := tiles.NewTilemapJSON("./assets/maps/basic_map.json")

	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()

	if err != nil {
		log.Fatal(err)
	}

	wallsAndFloors := tilemapJSON.GenColliders()

	playerImg, _, err := ebitenutil.NewImageFromFile("./assets/images/green_character_spritesheet.png")

	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(8, 4, 32)

	return &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   200.0,
				Y:   200.0,
			},
			Facing: entities.Left,
			State:  entities.Idle,
			Animations: map[entities.AnimationKey]*animation.Animation{
				{State: entities.Idle, Facing: entities.Left}:     animation.NewAnimation(0, 6, 1, 4.0),
				{State: entities.Idle, Facing: entities.Right}:    animation.NewAnimation(8, 14, 1, 4.0),
				{State: entities.Jumping, Facing: entities.Left}:  animation.NewAnimation(16, 21, 1, 2.0),
				{State: entities.Jumping, Facing: entities.Right}: animation.NewAnimation(24, 29, 1, 2.0),
				{State: entities.Running, Facing: entities.Left}:  animation.NewAnimation(22, 23, 1, 4.0),
				{State: entities.Running, Facing: entities.Right}: animation.NewAnimation(30, 31, 1, 4.0),
			},
		},
		tilemapJSON:       tilemapJSON,
		tilesets:          tilesets,
		tilemapImg:        tilemapImg,
		wallsAndFloors:    wallsAndFloors,
		playerSpriteSheet: playerSpriteSheet,
	}
}

func (g *Game) Update() error {
	g.player.Update(g.wallsAndFloors)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{0, 0, 0, 255})

	opts := ebiten.DrawImageOptions{}

	// loop over the layers
	for layerIndex, layer := range g.tilemapJSON.Layers {
		// loop over the tiles in the layer data
		for index, id := range layer.Data {

			if id == 0 {
				continue
			}

			// get the tile position of the tile
			x := index % layer.Width
			y := index / layer.Width

			// convert the tile position to pixel position
			x *= 32
			y *= 32

			img := g.tilesets[layerIndex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(img, &opts)

			// reset the opts for the next tile
			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Translate(g.player.X, g.player.Y)

	// PlayerSpeed = 5
	anim := g.player.ActiveAnimation()
	if anim != nil {
		frame := anim.Frame() // whatever your Animation exposes
		screen.DrawImage(
			g.player.Img.SubImage(
				g.playerSpriteSheet.Rect(frame),
			).(*ebiten.Image),
			&opts,
		)
	}

	for _, collider := range g.wallsAndFloors {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X),
			float32(collider.Min.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			true,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return constants.ScreenWidth, constants.ScreenHeight
}
