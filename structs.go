package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// VoxelDictionary, contains specific information about voxels.
type VoxelDictionary struct {
	Voxels []Voxel
}

// Voxel, contains information about a voxel.
type Voxel struct {
	Name    string
	Atlas   *ebiten.Image
	Texture [4]int
}

// VoxelPointer, has a reference to its voxel dictionary and an index.
type VoxelPointer struct {
	VoxelDictionary *VoxelDictionary
	Index           int
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
func (c *Chunk) GetVoxel(x, y, z int) VoxelPointer {
	return c.Voxels[x+y*c.Width+z*c.Width*c.Height]
}

// Make chunk from 3D array of voxels
func MakeChunk(voxels [][][]VoxelPointer) *Chunk {
	width, height, depth := len(voxels), len(voxels[0]), len(voxels[0][0])
	chunk := &Chunk{
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
	Chunks map[[2]int]*Chunk
}

// Return a Chunk from the world
func (w *World) GetChunk(x, y int) *Chunk {
	return w.Chunks[[2]int{x, y}]
}
