package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// json writable world metadata
type WorldMetadata struct {
	Seed     int64  `json:"seed"`
	SavePath string `json:"save_path"`
}

// json chunk
type DebugChunkJSON struct {
	// voxel name encoding map (for data compression)
	VoxelNamesShort map[string]int `json:"voxel_names_short"`
	VoxelNames      []int          `json:"voxel_names"`
	Width           int            `json:"width"`
	Height          int            `json:"height"`
	Depth           int            `json:"depth"`
}

// json player
type PlayerJSON struct {
	Position [3]float32 `json:"position"`
	Velocity [3]float32 `json:"velocity"`
}

/*
## LAYOUT
WORLD.json - world metadata & block table
PLAYER.json - player data
TERRAIN /
	CHUNK_X_Y.bin - RLE block data
	TAGS_X_Y.json - coordinate-based block data for blocks that have data tags (only if chunk has tagged blocks)
ENTITY /
	CHUNK_X_Y.json - coordinate-based entity data (only if chunk has entities)
*/

// fit metadata to world
func (world *World) ApplyMetadata(worldJSON WorldMetadata) {
	world.Seed = worldJSON.Seed
	world.SavePath = worldJSON.SavePath
}

// Convert a World to a WorldJSON
func (world *World) WorldToJSON() (worldJSON WorldMetadata) {
	worldJSON.Seed = world.Seed
	worldJSON.SavePath = world.SavePath

	return
}

// Convert a ChunkJSON to a Chunk
func (chunkJSON DebugChunkJSON) JSONToChunk() (chunk Chunk) {
	chunk.Depth = chunkJSON.Depth
	chunk.Width = chunkJSON.Width
	chunk.Height = chunkJSON.Height
	chunk.Voxels = make([]VoxelPointer, chunkJSON.Width*chunkJSON.Height*chunkJSON.Depth)
	for i := range chunk.Voxels {
		chunk.Voxels[i] = defaultVoxelDictionary.GetVoxelPointerTo(invertMap(chunkJSON.VoxelNamesShort)[chunkJSON.VoxelNames[i]])
	}

	return
}

// Convert a Chunk to a ChunkJSON DEBUG OBSOLETE
func (chunk Chunk) ChunkToJSON() (chunkJSON DebugChunkJSON) {
	chunkJSON.Depth = chunk.Depth
	chunkJSON.Width = chunk.Width
	chunkJSON.Height = chunk.Height
	chunkJSON.VoxelNames = make([]int, len(chunk.Voxels))
	chunkJSON.VoxelNamesShort = make(map[string]int)

	// voxel name encoding map (for data compression)

	// get all the unique voxel names
	allVoxelNames := make([]string, 0)
	for i := range chunk.Voxels {
		allVoxelNames = append(allVoxelNames, chunk.Voxels[i].GetVoxel().Name)
	}
	uniqueVoxelNames := uniqueItems(allVoxelNames)

	// create the encoding map
	for i, voxelName := range uniqueVoxelNames {
		chunkJSON.VoxelNamesShort[voxelName] = i
	}

	// encode the voxel names
	for i := range allVoxelNames {
		chunkJSON.VoxelNames[i] = chunkJSON.VoxelNamesShort[allVoxelNames[i]]
	}

	return
}

// convert player into json
func (player Player) ToJSON() (playerJSON PlayerJSON) {
	return PlayerJSON{
		[3]float32{player.Position.X, player.Position.Y, player.Position.Z},
		[3]float32{player.Velocity.X, player.Velocity.Y, player.Velocity.Z},
	}
}

// convert json to player
func (playerJSON PlayerJSON) ToPlayer() (player Player) {
	return Player{
		Position: Vec3{X: playerJSON.Position[0], Y: playerJSON.Position[1], Z: playerJSON.Position[2]},
		Velocity: Vec3{X: playerJSON.Velocity[0], Y: playerJSON.Velocity[1], Z: playerJSON.Velocity[2]},
		Drag:     Vec3{.9, .9, .9},
		Texture:  "Default",
	}
}

// create empty save, overwrites
func (game *Game) MakeEmptySave() (err error) {
	// ensure the save path exists
	err = os.MkdirAll(game.World.SavePath, 0755)
	if err != nil {
		log.Printf("ERROR: Failed to create save path: %v", err)
		return
	}
	os.MkdirAll(filepath.Join(game.World.SavePath, "terrain"), 0755)
	os.MkdirAll(filepath.Join(game.World.SavePath, "entity"), 0755) // why this code?

	return nil
}

// write non-chunk-data
func (game *Game) WriteData() (err error) {
	worldJSON := game.World.WorldToJSON()
	playerJSON := game.Player.ToJSON()

	// marshal the things to json
	worldJSONData, err := json.Marshal(worldJSON)
	if err != nil {
		log.Printf("ERROR: Failed to marshal world: %v", err)
		return
	}
	playerJSONData, err := json.Marshal(playerJSON)

	// write the json to a file
	err = os.WriteFile(filepath.Join(game.World.SavePath, "world.json"), worldJSONData, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to write world metadata: %v", err)
		return
	}
	err = os.WriteFile(filepath.Join(game.World.SavePath, "player.json"), playerJSONData, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to write player save: %v", err)
		return
	}

	return
}

// read non-chunk-data
// load a chunk from file
func (game *Game) LoadData() (err error) {
	// open the file
	var worldFile *os.File
	var playerFile *os.File
	worldFile, err = os.Open(filepath.Join(game.World.SavePath, "world.json"))
	if err != nil {
		return
	}
	playerFile, err = os.Open(filepath.Join(game.World.SavePath, "player.json"))
	if err != nil {
		return
	}

	// parse the text
	worldDecoder := json.NewDecoder(worldFile)
	defer worldFile.Close()
	playerDecoder := json.NewDecoder(playerFile)
	defer playerFile.Close()

	var worldMetadata WorldMetadata
	var playerJSON PlayerJSON
	err = worldDecoder.Decode(&worldMetadata)
	if err != nil {
		return err
	}
	err = playerDecoder.Decode(&playerJSON)
	if err != nil {
		return err
	}

	// convert
	player := playerJSON.ToPlayer()

	// load
	game.World.ApplyMetadata(worldMetadata)
	game.Player.Position = player.Position
	game.Player.Velocity = player.Velocity

	return nil
}

// make the filename for a chunk save
func chunkFileNameFromCoordinate(x, y int) string {
	return fmt.Sprintf("chunk%v_%v.json", x, y)
}

// read the coordinates for a chunk save file
func chunkCoordinateFromFileName(fileName string) (x, y int, err error) {
	coordinateString := strings.Replace(fileName, "chunk", "", 1)
	splitString := strings.Split(coordinateString, "_")

	x, err = strconv.Atoi(splitString[0])
	if err != nil {
		return
	}
	y, err = strconv.Atoi(splitString[1])
	if err != nil {
		return
	}

	return
}

// write a chunk to a file
func (world *World) WriteChunk(chunk Chunk, x, y int) (err error) {
	chunkJSON := chunk.ChunkToJSON()

	if pathExists(filepath.Join(world.SavePath, "world.json")) {
		// marshal the chunk to json
		jsonData, n_err := json.Marshal(chunkJSON)
		if n_err != nil {
			log.Printf("ERROR: Failed to marshal chunk: %v", err)
			return
		}

		// write the json to a file
		saveName := chunkFileNameFromCoordinate(x, y)
		err = os.WriteFile(filepath.Join(world.SavePath, "terrain", saveName), jsonData, 0644) // why this code?
		if err != nil {
			log.Printf("ERROR: Failed to write chunk: %v", err)
			return
		}
	} else {
		return fmt.Errorf("World does not exist!")
	}

	return nil
}

// load a chunk from file
func (world *World) LoadChunk(x, y int) (chunk Chunk, err error) {
	fileName := chunkFileNameFromCoordinate(x, y)
	path := filepath.Join(world.SavePath, "terrain", fileName)

	if pathExists(path) {
		// open the file
		var file *os.File
		file, err = os.Open(path)
		if err != nil {
			return Chunk{}, err
		}

		// parse the text
		decoder := json.NewDecoder(file)
		defer file.Close()

		var chunkJSON DebugChunkJSON
		err = decoder.Decode(&chunkJSON)
		if err != nil {
			return Chunk{}, err
		}

		// convert the ChunkJSON to a Chunk
		chunk = chunkJSON.JSONToChunk()

		return chunk, nil
	} else {
		return Chunk{}, fmt.Errorf("Chunk does not exist!") // empty chunk
	}
	return
}

// load game. does not load any chunks
func (game *Game) LoadGame(savePath string) (err error) {
	// load world metadata
	if pathExists(filepath.Join(savePath, "world.json")) {
		// open the world file
		var file *os.File
		file, err = os.Open(filepath.Join(savePath, "world.json"))
		if err != nil {
			return
		}

		// parse the text
		decoder := json.NewDecoder(file)
		defer file.Close()

		var worldJSON WorldMetadata
		err = decoder.Decode(&worldJSON)
		if err != nil {
			return
		}
		game.World.Initialize(worldJSON.Seed)
		game.World.ApplyMetadata(worldJSON)

		// load the rest of the data
		game.LoadData()
	} else {
		return fmt.Errorf("Game save does not exist!") // empty chunk
	}
	return
}

// check if there is a file at the given path
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// chunk exists
func (world *World) chunkExists(x, y int) bool {
	return pathExists(filepath.Join(world.SavePath, "terrain", chunkFileNameFromCoordinate(x, y)))
}

// save routine

// chunk loading go routine
func (game *Game) syncWorldWithDisk() {
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
					// unload and save
					for _, key := range chunksToUnload {
						// save
						chunk := game.World.Chunks[key]
						err := game.World.WriteChunk(chunk, key[0], key[1])
						if err != nil {
							panic(err)
						}
						// unload
						delete(game.World.Chunks, key)
					}

					// load and generate
					for _, key := range chunksToLoad {
						if game.World.chunkExists(key[0], key[1]) {
							chunk, err := game.World.LoadChunk(key[0], key[1])
							if err != nil {
								panic("Unable to load chunk.")
							}
							game.World.Chunks[key] = chunk
						} else {
							game.World.generateChunk(key, game.World.ChunkSize, game.World.ChunkSize, game.World.ChunkDepth, defaultVoxelDictionary)
						}
					}
				}

				// save all the random data
				game.WriteData()
			}
		}
	}
}
