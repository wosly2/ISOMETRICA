package main

import (
	"math"

	"math/rand"
)

// generate a procedurally generated chunk
func generateChunk(chunkWidth, chunkDepth, chunkHeight int, VDict VoxelDictionary, seed int64) Chunk {
	rand.Seed(seed)

	// create empty chunk
	chunkArray := make([][][]VoxelPointer, chunkWidth)
	for i := 0; i < chunkWidth; i++ {
		chunkArray[i] = make([][]VoxelPointer, chunkHeight)
		for j := 0; j < chunkHeight; j++ {
			chunkArray[i][j] = make([]VoxelPointer, chunkDepth)
		}
	}
	chunk := MakeChunk(chunkArray)

	// procedurally generate voxels
	for x := 0; x < chunkWidth; x++ {
		for y := 0; y < chunkHeight; y++ {
			for z := 0; z < chunkDepth; z++ {
				// set to air by default
				chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 0})

				// fill with water up to 5
				if z < 5 {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 2})
				}

				// fill with dirt up to int(math.Sin(float64(x)/4)*10+math.Sin(float64(y))*2)
				if z < int(math.Sin(float64(x)/4)*10+math.Sin(float64(y))*2) {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 5})
				}

				// fill with grass at int(math.Sin(float64(x)/4)*10+math.Sin(float64(y))*2)
				if z == int(math.Sin(float64(x)/4)*10+math.Sin(float64(y))*2) {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 1})
				}
			}
		}
	}

	// decorations

	// flowers
	for x := 0; x < chunkWidth; x++ {
		for y := 0; y < chunkHeight; y++ {
			if rand.Intn(10) == 0 {
				chunk.PlaceFlower(x, y)
			}
		}
	}

	return chunk
}

// May not work 100% of the time
func (chunk *Chunk) PlaceFlower(x, y int) (placed bool) {
	placed = false
	// move down until we hit a grass block that has have air above it
	for z := chunk.Depth - 1; z >= 0; z-- {
		if chunk.GetVoxel(x, y, z).Name == "Grass" {
			if chunk.GetVoxel(x, y+1, z).Name == "Air" {
				chunk.SetVoxel(x, y, z+1, VoxelPointer{VoxelDictionary: &defaultVoxelDictionary, Index: 8})
				placed = true
				break
			}
		}
	}
	return
}

// place a tree at the given position
func (chunk *Chunk) PlaceTree(x, y int) (placed bool) {
	placed = false
	// move down until we hit a grass block that has have air above it
	for z := chunk.Depth - 1; z >= 0; z-- {
		if chunk.GetVoxel(x, y, z).Name == "Grass" {
			// check if there is a 3x3x5 block of open air above centered above it
			for x2 := x - 1; x2 <= x+1; x2++ {
				for y2 := y - 1; y2 <= y+1; y2++ {
					for z2 := z - 1; z2 <= z+1; z2++ {
						if chunk.GetVoxel(x2, y2, z2).Name != "Air" {
							return false
						}
					}
				}
			}
			// place the leaves
			// place the trunk
		}
	}
	return
}

// generate a chunk that fits into the world via moving the perlin noise to the chunk location
func generateChunkBasedOnWorldPosition(x, y int, VDict VoxelDictionary, seed int64) Chunk {
	// TODO: Implement this instead of just returning a generated chunk
	return generateChunk(16, 16, 16, VDict, seed-int64(x)-int64(y))
}
