package main

import "github.com/hajimehoshi/ebiten/v2"

// listen to inputs
func runStateInput(game *Game) {
	// var inputs []string

	var playerSpeed float32 = .05

	// player movement
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		game.Player.Velocity.Y -= playerSpeed
		// inputs = append(inputs, "W")
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		game.Player.Velocity.Y += playerSpeed
		// inputs = append(inputs, "S")
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		game.Player.Velocity.X -= playerSpeed
		// inputs = append(inputs, "A")
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		game.Player.Velocity.X += playerSpeed
		// inputs = append(inputs, "D")
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		game.Player.Velocity.Z += playerSpeed
		// inputs = append(inputs, "Space")
	}
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		game.Player.Velocity.Z -= playerSpeed
		// inputs = append(inputs, "ShiftLeft")
	}

	// pause
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		game.GameState = GAMESTATE_MENU
	}

	// toggle depth shift
	if ebiten.IsKeyPressed(ebiten.KeyBackslash) {
		game.UsingDepthShift = true
		// inputs = append(inputs, "Backslash")
	} else {
		game.UsingDepthShift = false
		game.DepthShift = 0
	}

	// toggle debug mode
	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		game.DebugMode = !game.DebugMode
		ebiten.SetVsyncEnabled(!game.DebugMode)
		// inputs = append(inputs, "F3")
	}

	// rotate camera
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		game.Direction = NORTH
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		game.Direction = SOUTH
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		game.Direction = WEST
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		game.Direction = EAST
	}

	// if len(inputs) > 0 {
	// 	log.Printf("Inputs: %s\n", strings.Join(inputs, ", "))
	// }
}

func gameStateUpdateRun(game *Game) error {
	runStateInput(game)

	// depth shift
	if game.UsingDepthShift {
		if game.DepthShift > 5 {
			game.DepthShiftDirection = false
		} else if game.DepthShift < -5 {
			game.DepthShiftDirection = true
		}

		if game.DepthShiftDirection {
			game.DepthShift += .3
		} else {
			game.DepthShift -= .3
		}
	}

	// update player
	game.Player.Update(game.World)

	// get current chunk based on player position
	game.CurrentChunk[0] = int(game.Player.Position.X / float32(game.World.ChunkSize))
	game.CurrentChunk[1] = int(game.Player.Position.Y / float32(game.World.ChunkSize))

	// change camera position to have player in the center
	playerX, playerY := game.Player.getScreenPosition(game.DepthShift)
	game.Camera[0] = float32(game.ScreenX/2) - playerX
	game.Camera[1] = float32(game.ScreenY/2) - playerY

	return nil
}
