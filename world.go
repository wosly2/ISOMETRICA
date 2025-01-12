package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/aquilax/go-perlin"
	"github.com/hajimehoshi/ebiten/v2"
)

// VoxelDictionary, contains specific information about voxels.
type VoxelDictionary struct {
	Voxels               []Voxel
	Transparent          []string
	TransparentNoCulling []string
	Opaque               []string
}

// get a []string of voxels that are transparent
func (vDict *VoxelDictionary) GetTranparentNames() (transparentNames []string) {
	transparentNames = make([]string, len(vDict.Transparent)+len(vDict.TransparentNoCulling))
	for i := 0; i < len(transparentNames); i++ {
		if i < len(vDict.Transparent) {
			transparentNames[i] = vDict.Transparent[i]
		} else {
			transparentNames[i] = vDict.TransparentNoCulling[i-len(vDict.Transparent)]
		}
	}
	return
}

// return a pointer to a voxel in the dictionary
func (vDict *VoxelDictionary) GetVoxelPointerTo(name string) (pointer VoxelPointer) {
	pointer = VoxelPointer{&errorVoxelDictionary, 0}
	for i := 0; i < len(vDict.Voxels); i++ {
		if vDict.Voxels[i].Name == name {
			pointer = VoxelPointer{vDict, i}
			break
		}
	}
	return
}

// return a voxel address based on name
func (vDict *VoxelDictionary) GetVoxelNamed(name string) (voxel Voxel) {
	for i := 0; i < len(vDict.Voxels); i++ {
		if vDict.Voxels[i].Name == name {
			voxel = vDict.Voxels[i]
			break
		}
	}
	return
}

// Voxel, contains information about a voxel.
type Voxel struct {
	Name        string
	Atlas       *ebiten.Image
	TextureRect image.Rectangle
}

// VoxelPointer, has a reference to its voxel dictionary and an index.
// VoxelPointers and their VoxelDictionary seem to be on track to become deprecated by *Voxel.
// The only obvious benefit of VoxelPointer is the organization provided by the VoxelDicts.
// Their added complexity is starting to make code hard to read.
type VoxelPointer struct {
	VoxelDictionary *VoxelDictionary
	Index           int
}

// return the real Voxel from the VoxelPointer
func (pointer *VoxelPointer) GetVoxel() Voxel {
	return pointer.VoxelDictionary.Voxels[pointer.Index]
}

// Chunk, contains voxels.
// Voxels are stored in a 1D array.
type Chunk struct {
	Voxels []VoxelPointer
	Width  int
	Height int
	Depth  int
}

// Get voxel at x, y, z
func (c *Chunk) GetVoxel(x, y, z int) (voxel Voxel) {
	// check if voxel is out of bounds
	if x < 0 || y < 0 || z < 0 || x >= c.Width || y >= c.Height || z >= c.Depth {
		//log.Printf("Voxel {%d, %d, %d} out of bounds\n", x, y, z)
		return errorVoxelDictionary.GetVoxelNamed("Error")
	}
	var voxelPointer = c.Voxels[x+y*c.Width+z*c.Width*c.Height]
	voxel = voxelPointer.VoxelDictionary.Voxels[voxelPointer.Index]
	return voxel
}

// get voxel dictionary at x, y, z
func (c *Chunk) GetVoxelDictionary(x, y, z int) *VoxelDictionary {
	var voxelPointer = c.Voxels[x+y*c.Width+z*c.Width*c.Height]
	return voxelPointer.VoxelDictionary
}

// set voxel at x, y, z
func (c *Chunk) SetVoxel(x, y, z int, voxel VoxelPointer) (set bool) {
	// check if position is in bounds
	if x < 0 || y < 0 || z < 0 || x >= c.Width || y >= c.Height || z >= c.Depth {
		return false
	}
	// set voxel
	c.Voxels[x+y*c.Width+z*c.Width*c.Height] = voxel
	return true
}

// check if voxel is in bounds
func (c *Chunk) IsVoxelInBounds(x, y, z int) bool {
	return x >= 0 && y >= 0 && z >= 0 && x < c.Width && y < c.Height && z < c.Depth
}

// Make chunk from 3D array of voxels
func MakeChunk(voxels [][][]VoxelPointer) Chunk {
	width, height, depth := len(voxels), len(voxels[0]), len(voxels[0][0])
	chunk := Chunk{
		Voxels: make([]VoxelPointer, width*height*depth),
		Width:  width,
		Height: height,
		Depth:  depth,
	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			for z := 0; z < depth; z++ {
				chunk.Voxels[x+y*width+z*width*height] = voxels[x][y][z]
			}
		}
	}
	return chunk
}

// World, stores chunks in a map
type World struct {
	Chunks                 map[[2]int]Chunk
	Seed                   int64
	PerlinNoise            *perlin.Perlin
	WaterLevel             int
	SurfaceFeaturesBeginAt int
	ChunkSize              int
	ChunkDepth             int
	SavePath               string
	Initiated              bool
}

// Return a Chunk from the world
func (w *World) GetChunk(x, y int) (chunk Chunk, exists bool) {
	chunk, exists = w.Chunks[[2]int{x, y}]
	return
}

// return the voxel at x, y, z (global)
func (w *World) GetVoxel(x, y, z int) (voxel Voxel, exists bool) {
	chunk, exists := w.GetChunk(x/w.ChunkSize, y/w.ChunkSize)
	if exists {
		voxel = chunk.GetVoxel(x%w.ChunkSize, y%w.ChunkSize, z)
	}
	return
}

// json writable world
type WorldJSON struct {
	Chunks   map[string]ChunkJSON `json:"chunks"` // the string is the chunk position in the format "x,y"
	Seed     int64                `json:"seed"`
	SavePath string               `json:"save_path"`
}

// json chunk
type ChunkJSON struct {
	// voxel name encoding map (for data compression)
	VoxelNamesShort map[string]int `json:"voxel_names_short"`
	VoxelNames      []int          `json:"voxel_names"`
	Width           int            `json:"width"`
	Height          int            `json:"height"`
	Depth           int            `json:"depth"`
}

// world save structure:
// save file (name of world)
// // world.json (world json data)

// make ChunkJSON key
func ChunkJSONKey(key [2]int) string {
	return fmt.Sprintf("%d,%d", key[0], key[1])
}

// parse ChunkJSON key
func ParseChunkJSONKey(key string) (parsedKey [2]int, err error) {
	var x, y int
	_, err = fmt.Sscanf(key, "%d,%d", &x, &y)
	parsedKey = [2]int{x, y}
	return
}

// Convert a WorldJSON to a World
func (worldJSON *WorldJSON) JSONToWorld() (world World) {
	world = World{
		Chunks:     make(map[[2]int]Chunk),
		Seed:       worldJSON.Seed,
		ChunkSize:  32,
		ChunkDepth: 64,
		SavePath:   worldJSON.SavePath,
	}

	for key := range worldJSON.Chunks {
		parsedKey, _ := ParseChunkJSONKey(key)
		world.Chunks[parsedKey] = worldJSON.Chunks[key].JSONToChunk()
	}

	return
}

// Convert a World to a WorldJSON
func (world *World) WorldToJSON() (worldJSON WorldJSON) {
	worldJSON.Chunks = make(map[string]ChunkJSON)
	worldJSON.Seed = world.Seed
	worldJSON.SavePath = world.SavePath

	for key := range world.Chunks {
		worldJSON.Chunks[ChunkJSONKey(key)] = world.Chunks[key].ChunkToJSON()
	}

	return
}

// Convert a ChunkJSON to a Chunk
func (chunkJSON ChunkJSON) JSONToChunk() (chunk Chunk) {
	chunk.Depth = chunkJSON.Depth
	chunk.Width = chunkJSON.Width
	chunk.Height = chunkJSON.Height
	chunk.Voxels = make([]VoxelPointer, chunkJSON.Width*chunkJSON.Height*chunkJSON.Depth)
	for i := range chunk.Voxels {
		chunk.Voxels[i] = defaultVoxelDictionary.GetVoxelPointerTo(invertMap(chunkJSON.VoxelNamesShort)[chunkJSON.VoxelNames[i]])
	}

	return
}

// Convert a Chunk to a ChunkJSON
func (chunk Chunk) ChunkToJSON() (chunkJSON ChunkJSON) {
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

func readWorld(savePath string) (world World, err error) {
	// open the file
	file, err := os.Open(filepath.Join(savePath, "world.json"))
	if err != nil {
		log.Printf("ERROR: Failed to open world: %v", err)
		return
	}

	// parse the text
	decoder := json.NewDecoder(file)
	defer file.Close()

	var worldJSON WorldJSON
	err = decoder.Decode(&worldJSON)
	if err != nil {
		log.Printf("ERROR: Failed to parse world: %v", err)
		return
	}

	// convert the WorldJSON to a World
	world = worldJSON.JSONToWorld()

	// initialize the world
	world.Initalize(worldJSON.Seed)

	return
}

// write a world to a file. this DOES overwrite any existing world file.
// never use this to save a partially loaded world, as it will overwrite any currently unloaded chunks.
func (world World) WriteWorld(savePath string) (err error) {
	// ensure the save path exists
	err = os.MkdirAll(savePath, 0755)
	if err != nil {
		log.Printf("ERROR: Failed to create save path: %v", err)
		return
	}

	worldJSON := world.WorldToJSON()

	// marshal the world to json
	jsonData, err := json.Marshal(worldJSON)
	if err != nil {
		log.Printf("ERROR: Failed to marshal world: %v", err)
		return
	}

	// write the json to a file
	err = os.WriteFile(filepath.Join(savePath, "world.json"), jsonData, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to write world: %v", err)
		return
	}

	return
}

// check if there is a save file at the given path
func worldExists(savePath string) bool {
	_, err := os.Stat(filepath.Join(savePath, "world.json"))
	return !os.IsNotExist(err)
}
