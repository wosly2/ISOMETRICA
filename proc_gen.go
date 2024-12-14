package main

import (
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
)

func pseudoRandomTangent(x float64) float64 {
	return math.Tan(x*12.9898) - math.Floor(math.Tan(x*12.9898))
}

// initalize a world with things like random seed and perlin noise
func (world *World) Initalize(seed int64) {
	world.Seed = seed
	world.PerlinNoise = perlin.NewPerlin(
		float64(50+rand.Intn(20))/100,  // Persistence
		float64(100+rand.Intn(50))/100, // Lacunarity
		3,                              // Octaves
		world.Seed,                     // Seed
	)
	world.SurfaceFeaturesBeginAt = 10
	world.WaterLevel = 5 + world.SurfaceFeaturesBeginAt
}

// generate a procedurally generated chunk
func (world *World) generateChunk(position [2]int, chunkWidth, chunkHeight, chunkDepth int, VDict VoxelDictionary) {

	chunkArray := make([][][]VoxelPointer, chunkWidth) // Width (x)
	for x := 0; x < chunkWidth; x++ {
		chunkArray[x] = make([][]VoxelPointer, chunkHeight) // Height (y)
		for y := 0; y < chunkHeight; y++ {
			chunkArray[x][y] = make([]VoxelPointer, chunkDepth) // Depth (z)
			for z := 0; z < chunkDepth; z++ {
				chunkArray[x][y][z] = VoxelPointer{VoxelDictionary: &VDict, Index: 0}
			}
		}
	}
	chunk := MakeChunk(chunkArray)

	// procedurally generate voxels
	for x := 0; x < chunkWidth; x++ {
		for y := 0; y < chunkHeight; y++ {
			for z := 0; z < chunkDepth; z++ {

				// get the noise value at this position
				// log.Print("POSITION DEBUG :(")
				// log.Printf("Local position: %d, %d\n", x, y)
				// log.Printf("Chunk position: %d, %d\n", position[0], position[1])
				// log.Printf("Chunk size: %d, %d\n", chunkWidth, chunkHeight)
				// log.Printf("Global position: %d, %d\n", (position[0]*chunkWidth)+x, (position[1]*chunkHeight)+y)
				var scale float64 = .02
				noiseValue := world.PerlinNoise.Noise2D(float64(position[0]*chunkWidth+x)*scale, float64(position[1]*chunkHeight+y)*scale) * 6
				if noiseValue < float64(world.WaterLevel) {
					noiseValue *= 1.5
				}

				// set to air by default
				chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 0})

				// fill with water up to the water level
				if z <= world.WaterLevel {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 2})
				}

				// fill with dirt up noise value
				if z <= world.SurfaceFeaturesBeginAt+int(noiseValue)+2 {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 5})
				}

				// fill with grass
				if z == world.SurfaceFeaturesBeginAt+int(noiseValue)+2 && world.SurfaceFeaturesBeginAt+int(noiseValue)+2 >= world.WaterLevel {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 1})
				}

				// fill with stone up to a point
				if z <= world.SurfaceFeaturesBeginAt+int(noiseValue)-1 || z <= world.SurfaceFeaturesBeginAt-(2+int(noiseValue/10)) {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 4})
				}

				// ocean sand
				if z == world.SurfaceFeaturesBeginAt+2 && chunk.GetVoxel(x, y, z+1).Name == "Water" && chunk.GetVoxel(x, y, z).Name != "Water" {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 3})
				}
			}
		}
	}

	// decorations
	for x := 0; x < chunkWidth; x++ {
		for y := 0; y < chunkHeight; y++ {
			// grass
			if rand.Intn(2) == 0 {
				chunk.PlaceDecoration(x, y, defaultVoxelDictionary.GetVoxelPointerTo("Tall_Grass"), defaultVoxelDictionary.GetVoxelPointerTo("Grass"))
			}

			// flowers
			if rand.Intn(5) == 0 {
				chunk.PlaceDecoration(x, y, defaultVoxelDictionary.GetVoxelPointerTo("Flower"), defaultVoxelDictionary.GetVoxelPointerTo("Grass"))
			}

			// trees
			if rand.Intn(50) == 0 {
				chunk.PlaceTree(x, y)
			}
		}
	}

	world.Chunks[position] = chunk
}

// "Drop" a decoration onto a given voxel at a given (x, y) position.
// The algorithm starts at with z = chunk depth, and moves down until it finds that voxel below it the `placesOn` voxel.
// Then it places the decoration and returns true.
// It will only replace voxels that are Air.
func (chunk *Chunk) PlaceDecoration(x, y int, decoration, placesOn VoxelPointer) (placed bool) {
	placed = false
	// move down until we hit a grass block that has have air above it
	for z := chunk.Depth - 1; z >= 0; z-- {
		if chunk.GetVoxel(x, y, z-1).Name == placesOn.GetVoxel().Name {
			if chunk.GetVoxel(x, y, z).Name == "Air" {
				chunk.SetVoxel(x, y, z, decoration)
				placed = true
				break
			}
		}
	}
	return
}

// place a tree at the given position
func (chunk *Chunk) PlaceTree(x, y int) (placed bool) {
	// move down until we hit a grass block that has have air above it
	for z := chunk.Depth - 1; z >= 0; z-- {
		if chunk.GetVoxel(x, y, z).Name == "Grass" {
			// check if there is air above it
			if chunk.GetVoxel(x, y, z+1).Name != "Air" {
				return false
			}

			// leaves motherfucka
			for x2 := x - 1; x2 < x+2; x2++ {
				for y2 := y - 1; y2 < y+2; y2++ {
					chunk.SetVoxel(x2, y2, z+4, defaultVoxelDictionary.GetVoxelPointerTo("Leaves"))
					if math.Abs(float64(x2-x)) != math.Abs(float64(y2-y)) {
						chunk.SetVoxel(x2, y2, z+5, defaultVoxelDictionary.GetVoxelPointerTo("Leaves"))
					}
				}
			}
			chunk.SetVoxel(x, y, z+5, defaultVoxelDictionary.GetVoxelPointerTo("Leaves")) // top middle leaf

			// place the trunk
			for z2 := z + 1; z2 <= z+4; z2++ {
				chunk.SetVoxel(x, y, z2, defaultVoxelDictionary.GetVoxelPointerTo("Wood"))
			}
			return true
		}
	}
	return false
}
