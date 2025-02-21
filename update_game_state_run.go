package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func gameStateUpdateRun(game *Game) error {
	// INPUT

	var inputs []string

	var playerSpeed float32 = .05

	// player movement
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		game.Player.Velocity.Y -= playerSpeed
		inputs = append(inputs, "W")
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		game.Player.Velocity.Y += playerSpeed
		inputs = append(inputs, "S")
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		game.Player.Velocity.X -= playerSpeed
		inputs = append(inputs, "A")
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		game.Player.Velocity.X += playerSpeed
		inputs = append(inputs, "D")
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		game.Player.Velocity.Z += playerSpeed
		inputs = append(inputs, "Space")
	}
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		game.Player.Velocity.Z -= playerSpeed
		inputs = append(inputs, "ShiftLeft")
	}

	// pause
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		game.GameState = GAMESTATE_MENU
	}

	// toggle depth shift
	if ebiten.IsKeyPressed(ebiten.KeyBackslash) {
		game.UsingDepthShift = true
		inputs = append(inputs, "Backslash")
	} else {
		game.UsingDepthShift = false
		game.DepthShift = 0
	}

	// reload world
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		game.World.Initalize(rand.Int63n(1000000))
		// regenerate the chunk and its neighbors
		for x := -1; x < 2; x++ {
			for y := -1; y < 2; y++ {
				game.World.generateChunk([2]int{x + game.CurrentChunk[0], y + game.CurrentChunk[1]}, game.World.ChunkSize, game.World.ChunkSize, game.World.ChunkDepth, defaultVoxelDictionary)
			}
		}
		inputs = append(inputs, "R")
	}

	// toggle debug mode
	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		game.DebugMode = !game.DebugMode
		ebiten.SetVsyncEnabled(!game.DebugMode)
		inputs = append(inputs, "F3")
	}

	// if len(inputs) > 0 {
	// 	log.Printf("Inputs: %s\n", strings.Join(inputs, ", "))
	// }

	// END OF INPUT

	// depth shift
	if game.UsingDepthShift {
		if game.DepthShift > 1 {
			game.DepthShiftDirection = false
		} else if game.DepthShift < -1 {
			game.DepthShiftDirection = true
		}

		if game.DepthShiftDirection {
			game.DepthShift += .02
		} else {
			game.DepthShift -= .02
		}
	}

	// update player
	game.Player.Update(game.World)

	// get current chunk based on player position
	game.CurrentChunk[0] = int(game.Player.Position.X / float32(game.World.ChunkSize))
	game.CurrentChunk[1] = int(game.Player.Position.Y / float32(game.World.ChunkSize))

	// change camera position to have player in the center
	screenX, screenY := game.Player.getScreenPosition(game.DepthShift)
	game.Camera[0] = -screenX
	game.Camera[1] = -screenY

	return nil
}
