package main

import (
	"log"
	"time"
)

func (game *Game) Update() error {
	// initiate
	if !game.HasInitiatedUpdate {
		game.HasInitiatedUpdate = true

		go game.loadChunks()
	}

	var err error

	// game state
	switch game.GameState {
	case GAMESTATE_GAME:
		err = gameStateUpdateRun(game)
	case GAMESTATE_TITLE:
		err = gameStateUpdateTitle(game)
	case GAMESTATE_MENU:
		err = gameStateUpdateMenu(game)
	}

	return err
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
