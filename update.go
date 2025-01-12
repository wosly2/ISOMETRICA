package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func (game *Game) Update() error {
	if !game.HasInitiatedUpdate {
		game.HasInitiatedUpdate = true
		go game.loadChunks()
	}

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

// chunk loading go routine
func (game *Game) loadChunks() {
	// timer
	start := time.Now()
	first := true

	for {
		if time.Since(start).Seconds() >= IOtimeInterval || first {
			if first {
				first = false
			}
			start = time.Now()

			if true {
				// save and unload any chunks that are out of range, and attempt to load any chunks that are in range
				chunksToUnload := make([][2]int, 0)
				chunksToLoad := make([][2]int, 0)

				// get out of range chunks
				for key := range game.World.Chunks {
					if absi(key[0]-game.CurrentChunk[0]) > chunkLoadDistance || absi(key[1]-game.CurrentChunk[1]) > chunkLoadDistance {
						chunksToUnload = append(chunksToUnload, key)
					}
				}

				// get in range chunks (that aren't already loaded)
				for x := game.CurrentChunk[0] - chunkLoadDistance; x <= game.CurrentChunk[0]+chunkLoadDistance; x++ {
					for y := game.CurrentChunk[1] - chunkLoadDistance; y <= game.CurrentChunk[1]+chunkLoadDistance; y++ {
						if _, exists := game.World.Chunks[[2]int{x, y}]; !exists {
							chunksToLoad = append(chunksToLoad, [2]int{x, y})
						}
					}
				}

				if len(chunksToUnload) > 0 || len(chunksToLoad) > 0 {
					// we will need the whole world on disk so we can load to it and from it
					var diskWorld World
					if worldExists(game.World.SavePath) {
						var err error
						diskWorld, err = readWorld(game.World.SavePath)
						if err != nil {
							log.Fatalf("Failed to load world from disk: %v", err)
						}
					} else {
						diskWorld = game.World
					}

					// unload chunks
					if len(chunksToUnload) > 0 {
						for _, key := range chunksToUnload {
							diskWorld.Chunks[key] = game.World.Chunks[key]
							delete(game.World.Chunks, key)
						}

						// write the world to disk
						err := diskWorld.WriteWorld(game.World.SavePath)
						if err != nil {
							log.Printf("ERROR: Failed to write world to disk: %v", err)
						}
					}

					// load chunks
					if len(chunksToLoad) > 0 {
						for _, key := range chunksToLoad {
							chunk, exists := diskWorld.Chunks[key]
							if exists {
								game.World.Chunks[key] = chunk
							} else {
								game.World.generateChunk(key, game.World.ChunkSize, game.World.ChunkSize, game.World.ChunkDepth, defaultVoxelDictionary)
							}
						}
					}
				}
			}
		}
	}
}
