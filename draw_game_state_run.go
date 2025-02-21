package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func gameStateDrawRun(game *Game, screen *ebiten.Image) error {
	// fill background
	game.Framebuffer.Fill(color.RGBA{0, 0, 88, 255})

	// render the chunks
	var blocksRendered int
	for x := -2; x <= 2; x++ {
		for y := -2; y <= 2; y++ {
			// check if chunk is on screen
			if !ChunkContainingGlobalPointVisibleInViewport((game.CurrentChunk[0]+x)*32, (game.CurrentChunk[1]+y)*32, 0, game.Camera[0], game.Camera[1], game.DepthShift, game.ScreenX, game.ScreenY) {
				continue
			}

			chunk, exists := game.World.Chunks[[2]int{x + game.CurrentChunk[0], y + game.CurrentChunk[1]}]
			if exists {
				// get screen position of the top voxel in the chunk
				screenX, screenY := getScreenPosition(
					(x+game.CurrentChunk[0])*game.World.ChunkSize,
					(y+game.CurrentChunk[1])*game.World.ChunkSize,
					0,
					game.Camera[0],
					game.Camera[1],
					game.DepthShift,
				)
				// render
				blocksRendered += chunk.Render(game.Framebuffer, float32(screenX), float32(screenY), game.DepthShift, game, true)
			}
		}
	}

	if game.DebugMode {
		// draw the global world origin
		originX, originY := getScreenPosition(0, 0, 0, game.Camera[0], game.Camera[1], game.DepthShift)
		// x
		offsetX, offsetY := getScreenPosition(10, 0, 0, game.Camera[0], game.Camera[1], game.DepthShift)
		ebitenutil.DrawLine(
			game.Framebuffer,
			float64(originX),
			float64(originY),
			float64(offsetX),
			float64(offsetY),
			color.RGBA{255, 0, 0, 255},
		)
		// y
		offsetX, offsetY = getScreenPosition(0, 10, 0, game.Camera[0], game.Camera[1], game.DepthShift)
		ebitenutil.DrawLine(
			game.Framebuffer,
			float64(originX),
			float64(originY),
			float64(offsetX),
			float64(offsetY),
			color.RGBA{0, 255, 0, 255},
		)
		// z
		offsetX, offsetY = getScreenPosition(0, 0, 10, game.Camera[0], game.Camera[1], game.DepthShift)
		ebitenutil.DrawLine(
			game.Framebuffer,
			float64(originX),
			float64(originY),
			float64(offsetX),
			float64(offsetY),
			color.RGBA{0, 0, 255, 255},
		)
	}

	// draw the player
	game.Player.Render(game.Framebuffer, game.DepthShift)

	// gui/text
	drawString(game.Framebuffer, "ISOMETRICA Infdev", int(0), int(0), true)
	if game.DebugMode {
		drawString(game.Framebuffer, fmt.Sprintf("Player Position: %f, %f, %f", game.Player.Position.X, game.Player.Position.Y, game.Player.Position.Z), 0, 22, true)
		drawString(game.Framebuffer, fmt.Sprintf("Local position: %f, %f, %f", game.Player.Position.X-float32(game.CurrentChunk[0]*game.World.ChunkSize), game.Player.Position.Y-float32(game.CurrentChunk[1]*game.World.ChunkSize), game.Player.Position.Z), 0, 34, true)
		drawString(game.Framebuffer, fmt.Sprintf("Focused Chunk: %d, %d", game.CurrentChunk[0], game.CurrentChunk[1]), 0, 46, true)
		drawString(game.Framebuffer, fmt.Sprintf("FPS: %f TPS: %f", game.ActualFPS, ebiten.ActualTPS()), 0, 58, true)
		drawString(game.Framebuffer, fmt.Sprintf("Blocks Rendered: %d", blocksRendered), 0, 70, true)
		drawString(game.Framebuffer, fmt.Sprintf("Chunks Loaded: %d, World Byte Size: %d", len(game.World.Chunks), MapSize(game.World.Chunks)), 0, 82, true)
		/* velocity */ drawString(game.Framebuffer, fmt.Sprintf("Velocity: %f, %f, %f", game.Player.Velocity.X, game.Player.Velocity.Y, game.Player.Velocity.Z), 0, 94, true)
	} else {
		// drawString(game.Framebuffer, fmt.Sprintf("%f, %f, %f", game.Player.Position.X, game.Player.Position.Y, game.Player.Position.Z), 0, 22, true)
	}
	// draw framebuffer to screen
	screen.DrawImage(game.Framebuffer, nil)

	// frame control
	game.Frames++
	if time.Since(game.SecondTimer).Seconds() >= 1 {
		game.ActualFPS = float32(game.Frames)
		game.Frames = 0
		game.SecondTimer = time.Now()
	}

	return nil
}
