package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/Ethea2/feedball/constants"
	"github.com/Ethea2/feedball/entities"
	"github.com/Ethea2/feedball/tiles"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	player      *entities.Player
	tilemapJSON *tiles.TilemapJSON
	tilesets    []tiles.Tileset
	tilemapImg  *ebiten.Image
	colliders   []image.Rectangle
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

	colliders := tilemapJSON.GenColliders()

	playerImg, _, err := ebitenutil.NewImageFromFile("./assets/images/green_character_spritesheet.png")

	if err != nil {
		log.Fatal(err)
	}

	return &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   200.0,
				Y:   200.0,
			},
		},
		tilemapJSON: tilemapJSON,
		tilesets:    tilesets,
		tilemapImg:  tilemapImg,
		colliders:   colliders,
	}
}

func (g *Game) Update() error {
	g.player.Dx = 0.0
	g.player.Dy = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Dx = -constants.PlayerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Dx = constants.PlayerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Jump()
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Dy = constants.PlayerSpeed
	}

	if g.player.State == entities.Jumping {
		g.player.JumpTimer -= 0.02
		g.player.Dy -= constants.PlayerSpeed
		if g.player.JumpTimer <= 0 {
			g.player.State = entities.Down
			g.player.JumpTimer = 0
		}
	}

	fmt.Println(g.player.State)

	g.player.X += g.player.Dx

	CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy

	if g.player.State == entities.Down {
		g.player.Dy += constants.PlayerSpeed

		g.player.Y += g.player.Dy
	}

	CheckCollisionVertical(g.player.Sprite, g.colliders)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

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

	// draw the player
	// PlayerSpeed = 5
	screen.DrawImage(
		// grab a subimage of the spritesheet
		g.player.Img.SubImage(
			image.Rect(0, 0, 32, 32),
		).(*ebiten.Image),
		&opts,
	)

	for _, collider := range g.colliders {
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
