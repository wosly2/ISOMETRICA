package main

import (
	"image"

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
func (vDict *VoxelDictionary) GetTransparentNames() (transparentNames []string) {
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
