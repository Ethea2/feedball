# FeedBall 🏈

A 2D platformer game built with Go and [Ebitengine](https://ebitengine.org/).

## Overview

FeedBall is a tile-based 2D platformer where you control a character across a map loaded from a [Tiled](https://www.mapeditor.org/) JSON file. The game features tile rendering, collision detection, and basic player physics including jumping and gravity.

## Features

- Tile-based map rendering from Tiled `.json` map files
- Support for both uniform (spritesheet) and dynamic (individual image) tilesets
- Horizontal and vertical collision detection
- Player jumping with a timer-based arc
- Gravity simulation
- Debug collider outlines rendered in-game

## Project Structure

```
feedball/
├── main.go                  # Entry point, collision helpers
├── game.go                  # Game loop (Update, Draw, Layout)
├── constants/
│   └── constants.go         # Screen size, tile size, player speed
├── entities/
│   ├── sprite.go            # Base sprite struct (position, velocity, image)
│   └── player.go            # Player struct, states, jump logic
├── tiles/
│   ├── tilemap.go           # Tiled JSON parsing, collider generation
│   └── tileset.go           # Tileset loading (uniform & dynamic)
└── assets/
    ├── images/              # Spritesheets and terrain images
    └── maps/                # Tiled map files (.json, .tsj, .tmj)
```

## Prerequisites

- [Go](https://go.dev/) 1.21+
- Ebitengine dependencies (handled automatically via `go mod`)

On Linux, Ebitengine requires some system libraries. Install them with:

```bash
sudo apt install libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
```

## Getting Started

```bash
# Clone the repository
git clone https://github.com/Ethea2/feedball.git
cd feedball

# Download dependencies
go mod tidy

# Run the game
go run .
```

## Controls

| Key | Action |
|-----|--------|
| ← Left Arrow | Move left |
| → Right Arrow | Move right |
| ↑ Up Arrow | Jump |
| ↓ Down Arrow | Move down |

## Map Format

Maps are created with [Tiled Map Editor](https://www.mapeditor.org/) and exported as JSON. Place map files in `assets/maps/` and tileset files (`.tsj`) alongside them. The tilemap loader resolves tileset paths relative to the `assets/maps/` directory.

Tiles with a non-zero ID are treated as solid and generate collision rectangles automatically.

## Dependencies

- [Ebitengine v2](https://github.com/hajimehoshi/ebiten) — 2D game engine for Go
