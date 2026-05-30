package main

import (
	"image"
	"image/color"
	"log"

	"github.com/Ethea2/feedball/animation"
	"github.com/Ethea2/feedball/entities"
	"github.com/Ethea2/feedball/shared"
	"github.com/Ethea2/feedball/spritesheet"
	"github.com/Ethea2/feedball/tiles"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	ball              *entities.Ball
	tilemapJSON       *tiles.TilemapJSON
	tilesets          []tiles.Tileset
	tilemapImg        *ebiten.Image
	wallsAndFloors    []image.Rectangle
	releaseFrames     map[shared.Direction][2]int
	prevShootPressed  bool
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

	playerImg, _, err := ebitenutil.NewImageFromFile(
		"./assets/images/P1Sheet.png",
	)

	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(12, 8, 32)

	ballImg, _, err := ebitenutil.NewImageFromFile("./assets/images/Ball.png")

	if err != nil {
		log.Fatal(err)
	}

	ballSpawnX, ballSpawnY := tiles.TileToPixel(20, 10)

	return &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img:          playerImg,
				X:            200.0,
				Y:            200.0,
				LeftXOffset:  10,
				RightXOffset: -10,
				DownYOffset:  0,
			},
			Facing: entities.Left,
			State:  entities.Idle,
			Animations: map[entities.AnimationKey]*animation.Animation{
				{State: entities.VisualIdle, Facing: entities.Left}: animation.NewAnimation(
					0,
					4,
					1,
					4.0,
				),
				{State: entities.VisualIdle, Facing: entities.Right}: animation.NewAnimation(
					12,
					16,
					1,
					4.0,
				),
				{State: entities.VisualJumping, Facing: entities.Left}: animation.NewAnimation(
					8,
					10,
					1,
					2.0,
				),
				{State: entities.VisualJumping, Facing: entities.Right}: animation.NewAnimation(
					21,
					23,
					1,
					2.0,
				),
				{State: entities.VisualRunning, Facing: entities.Left}: animation.NewAnimation(
					6,
					7,
					1,
					4.0,
				),
				{State: entities.VisualRunning, Facing: entities.Right}: animation.NewAnimation(
					18,
					19,
					1,
					4.0,
				),
				{State: entities.VisualHoldingBallIdle, Facing: entities.Left}: animation.NewAnimation(
					24,
					28,
					1,
					4.0,
				),
				{State: entities.VisualHoldingBallIdle, Facing: entities.Right}: animation.NewAnimation(
					36,
					40,
					1,
					4.0,
				),
				{State: entities.VisualHoldingBallJumping, Facing: entities.Left}: animation.NewAnimation(
					32,
					35,
					1,
					2.0,
				),
				{State: entities.VisualHoldingBallJumping, Facing: entities.Right}: animation.NewAnimation(
					42,
					47,
					1,
					4.0,
				),
				{State: entities.VisualHoldingBallRunning, Facing: entities.Left}: animation.NewAnimation(
					30,
					31,
					1,
					4.0,
				),
				{State: entities.VisualHoldingBallRunning, Facing: entities.Right}: animation.NewAnimation(
					42,
					43,
					1,
					4.0,
				),
				{State: entities.VisualThrowingWindup, Facing: entities.Left}: animation.NewAnimation(
					48,
					53,
					1,
					shared.ThrowFrameSpeed,
				),
				{State: entities.VisualThrowingWindup, Facing: entities.Right}: animation.NewAnimation(
					60,
					65,
					1,
					shared.ThrowFrameSpeed,
				),
			},
			IsAffectedByGravity: true,
		},
		ball: &entities.Ball{
			Sprite: &entities.Sprite{
				Img:          ballImg,
				X:            float64(ballSpawnX),
				Y:            float64(ballSpawnY),
				RightXOffset: -32,
				DownYOffset:  -32,
			},
			State: entities.Initial,
			Speed: shared.BallInitialSpeed,
		},
		tilemapJSON:       tilemapJSON,
		tilesets:          tilesets,
		tilemapImg:        tilemapImg,
		wallsAndFloors:    wallsAndFloors,
		playerSpriteSheet: playerSpriteSheet,
		releaseFrames: map[shared.Direction][2]int{
			shared.NorthEast: {84, 87},
			shared.East:      {76, 79},
			shared.SouthEast: {88, 91},
			shared.NorthWest: {54, 57},
			shared.West:      {66, 69},
			shared.SouthWest: {72, 75},
		},
	}
}

func (g *Game) Update() error {
	playerInput := entities.PlayerInput{
		Jump:  ebiten.KeyUp,
		Shoot: ebiten.KeyEnter,
		Left:  ebiten.KeyLeft,
		Right: ebiten.KeyRight,
		Down:  ebiten.KeyDown,
	}
	// Ball pickup
	if (g.ball.State == entities.Initial || g.ball.State == entities.BallInPlay) && g.ball.ThrownBy != g.player && g.player.Sprite.Bounds().Overlaps(g.ball.Sprite.Bounds()) {
		g.player.HoldingBall = true
		g.ball.State = entities.PlayerHeld
	}

	// Throw initiation
	shootPressed := ebiten.IsKeyPressed(playerInput.Shoot)
	if g.player.HoldingBall && !g.player.Throwing && shootPressed && !g.prevShootPressed {
		if g.player.Facing == entities.Left {
			g.player.StartThrow(48, 53)
		} else {
			g.player.StartThrow(60, 65)
		}
	}
	g.prevShootPressed = shootPressed

	// Throw completion
	if g.player.Throwing && g.player.ThrowAnimation != nil && g.player.ThrowAnimation.IsDone() {
		dir := g.player.FinishThrow()
		g.ball.Direction = dir
		g.ball.Speed = shared.BallInitialSpeed
		g.ball.State = entities.BallInPlay
		g.ball.X = g.player.X
		g.ball.Y = g.player.Y
		g.ball.ThrownBy = g.player
	}

	// Ball movement
	g.ball.Update()

	g.player.Update(
		g.wallsAndFloors,
		playerInput,
		g.releaseFrames,
	)

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

	opts.GeoM.Scale(shared.PlayerBallScale, shared.PlayerBallScale)
	opts.GeoM.Translate(g.player.X, g.player.Y)

	if g.player.Throwing && g.player.ThrowAnimation != nil {
		frame := g.player.ThrowAnimation.Frame()
		screen.DrawImage(
			g.player.Img.SubImage(
				g.playerSpriteSheet.Rect(frame),
			).(*ebiten.Image),
			&opts,
		)
	} else {
		anim := g.player.ActiveAnimation()
		if anim != nil {
			frame := anim.Frame()
			screen.DrawImage(
				g.player.Img.SubImage(
					g.playerSpriteSheet.Rect(frame),
				).(*ebiten.Image),
				&opts,
			)
		}
	}

	opts.GeoM.Reset()

	if !(g.ball.State == entities.PlayerHeld) {
		opts.GeoM.Scale(shared.PlayerBallScale, shared.PlayerBallScale)
		opts.GeoM.Translate(g.ball.X, g.ball.Y)
		screen.DrawImage(
			g.ball.Img,
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
	playerBounds := g.player.Bounds()
	vector.StrokeRect(
		screen,
		float32(playerBounds.Min.X),
		float32(playerBounds.Min.Y),
		float32(playerBounds.Dx()),
		float32(playerBounds.Dy()),
		1.0,
		color.RGBA{0, 255, 0, 255}, // green to distinguish from wall colliders
		true,
	)

	if !(g.ball.State == entities.PlayerHeld) {
		ballBounds := g.ball.Bounds()
		vector.StrokeRect(
			screen,
			float32(ballBounds.Min.X),
			float32(ballBounds.Min.Y),
			float32(ballBounds.Dx()),
			float32(ballBounds.Dy()),
			1.0,
			color.RGBA{0, 255, 0, 255}, // green to distinguish from wall colliders
			true,
		)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return shared.ScreenWidth, shared.ScreenHeight
}
