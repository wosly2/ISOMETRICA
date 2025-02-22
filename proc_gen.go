package main

import (
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
)

// VORONOI BIOME SYSTEM
// based on this POC: https://github.com/wosly2/py-chunked-voronoi/blob/main/main.py (my own code)

type Biome struct {
	Name string
}

var biomes []Biome = []Biome{
	{"Plains"},
	{"Snowy"},
	{"Forest"},
	//Biome{"Mountains"},
	{"Desert"},
}

type VoronoiPoint struct {
	ChunkX, ChunkY int
	X, Y           int
	Biome          *Biome
}

var VoronoiCache map[[2]int]VoronoiPoint

// gets the Voronoi point located in an arbitrary chunk
func getVoronoiPoint(chunkX, chunkY int) VoronoiPoint {
	// check if we have a cached point
	if VoronoiCache != nil {
		if cachedPoint, ok := VoronoiCache[[2]int{chunkX, chunkY}]; ok {
			return cachedPoint
		}
	}

	// seed the rng on chunk pos
	vrand := rand.New(rand.NewSource(int64(chunkX + chunkY)))

	point := VoronoiPoint{
		ChunkX: chunkX,
		ChunkY: chunkY,
		X:      vrand.Intn(32),
		Y:      vrand.Intn(32),
		Biome:  &biomes[vrand.Intn(len(biomes))],
	}

	// cache the point
	if VoronoiCache == nil {
		VoronoiCache = make(map[[2]int]VoronoiPoint)
	}
	VoronoiCache[[2]int{chunkX, chunkY}] = point

	// return seeded point
	return point
}

// gets the nearest Voronoi point at an arbitrary position
func getNearestVoronoiPoint(chunkX, chunkY, localX, localY int) VoronoiPoint {
	// convert x, y to global positions
	globalX := chunkX*32 + localX
	globalY := chunkY*32 + localY

	// get a list of vpoints in moore neighborhood of chunks
	var vPoints []VoronoiPoint
	for x := chunkX - 1; x <= chunkX+1; x++ {
		for y := chunkY - 1; y <= chunkY+1; y++ {
			vPoints = append(vPoints, getVoronoiPoint(x, y))
		}
	}

	// get the closest vpoint
	minDistance := math.MaxFloat64
	var closestVoronoiPoint VoronoiPoint
	for _, vPoint := range vPoints {
		// euclidean distance
		distance := math.Sqrt(math.Pow(float64(globalX-vPoint.X), 2) + math.Pow(float64(globalY-vPoint.Y), 2))
		// update closest
		if distance < minDistance {
			minDistance = distance
			closestVoronoiPoint = vPoint
		}
	}

	return closestVoronoiPoint
}

// get the biome at a given position
func getBiome(chunkX, chunkY, localX, localY int) *Biome {
	return getNearestVoronoiPoint(chunkX, chunkY, localX, localY).Biome
}

// idk why this is here
// func pseudoRandomTangent(x float64) float64 {
// 	return math.Tan(x*12.9898) - math.Floor(math.Tan(x*12.9898))
// }

// initalize a world with things like random seed and perlin noise
func (world *World) Initalize(seed int64) {
	world.Seed = seed
	world.PerlinNoise = perlin.NewPerlin(
		float64(50+rand.Intn(20))/100, // Persistence
		float64(rand.Intn(50))/100,    // Lacunarity
		3,                             // Octaves
		world.Seed,                    // Seed
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
				var scale float64 = .02
				noiseValue := world.PerlinNoise.Noise2D(float64(position[0]*chunkWidth+x)*scale, float64(position[1]*chunkHeight+y)*scale) * 10
				if z < world.WaterLevel { // makes underwater topography steeper
					noiseValue *= 2
				}
				if z > 30 { // "mountains"
					noiseValue *= 2
				}

				// get the biome at this position
				biome := getBiome(position[0], position[1], x, y)

				// set to air by default
				chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 0})

				// fill with water up to the water level
				if z <= world.WaterLevel {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 2})
				}

				// fill with dirt/sand up noise value
				if z <= world.SurfaceFeaturesBeginAt+int(noiseValue)+2 {
					var dirtBlock string

					if biome.Name == "Snowy" || biome.Name == "Plains" || biome.Name == "Forest" {
						dirtBlock = "Dirt"
					} else if biome.Name == "Mountains" {
						dirtBlock = "Stone"
					} else if biome.Name == "Desert" {
						dirtBlock = "Sand"
					}

					chunk.SetVoxel(x, y, z, defaultVoxelDictionary.GetVoxelPointerTo(dirtBlock))
				}

				// fill with grass
				if z == world.SurfaceFeaturesBeginAt+int(noiseValue)+2 && world.SurfaceFeaturesBeginAt+int(noiseValue)+2 >= world.WaterLevel {
					var grassBlock string

					if biome.Name == "Snowy" {
						grassBlock = "Snowy_Grass"
					} else if biome.Name == "Mountains" {
						grassBlock = "Stone"
					} else if biome.Name == "Desert" {
						grassBlock = "Sand"
					} else {
						grassBlock = "Grass"
					}

					chunk.SetVoxel(x, y, z, defaultVoxelDictionary.GetVoxelPointerTo(grassBlock))
				}

				// fill with stone up to a point
				if z <= world.SurfaceFeaturesBeginAt+int(noiseValue)-1 || z <= world.SurfaceFeaturesBeginAt-(2+int(noiseValue/10)) {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 4})
				}

				// ocean sand
				if z == world.SurfaceFeaturesBeginAt+2 && chunk.GetVoxel(x, y, z+1).Name == "Water" && chunk.GetVoxel(x, y, z).Name != "Water" {
					chunk.SetVoxel(x, y, z, VoxelPointer{VoxelDictionary: &VDict, Index: 3})
				}
				if chunk.GetVoxel(x, y, z).Name == "Dirt" && chunk.GetVoxel(x, y, z+1).Name == "Water" {
					chunk.SetVoxel(x, y, z, defaultVoxelDictionary.GetVoxelPointerTo("Sand"))
				}
			}
		}
	}

	// decorations
	for x := 0; x < chunkWidth; x++ {
		for y := 0; y < chunkHeight; y++ {
			biome := getBiome(position[0], position[1], x, y)

			var grassBlock string
			var flowerBlock string
			var grassDecoBlock string

			if biome.Name == "Snowy" {
				grassBlock = "Snowy_Grass"
				flowerBlock = "Snowy_Flower"
				grassDecoBlock = "Snowy_Tall_Grass"
			} else if biome.Name == "Desert" {
				grassBlock = "Sand"
			} else if biome.Name == "Mountains" {
				grassBlock = "Stone"
			} else {
				grassBlock = "Grass"
				flowerBlock = "Flower"
				grassDecoBlock = "Tall_Grass"
			}

			// grass
			if rand.Intn(2) == 0 && biome.Name != "Desert" && biome.Name != "Mountains" {
				chunk.PlaceDecoration(x, y, defaultVoxelDictionary.GetVoxelPointerTo(grassDecoBlock), defaultVoxelDictionary.GetVoxelPointerTo(grassBlock))
			}

			// flowers
			if rand.Intn(5) == 0 && biome.Name != "Desert" && biome.Name != "Mountains" {
				chunk.PlaceDecoration(x, y, defaultVoxelDictionary.GetVoxelPointerTo(flowerBlock), defaultVoxelDictionary.GetVoxelPointerTo(grassBlock))
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
		if chunk.GetVoxel(x, y, z).Name == "Grass" || chunk.GetVoxel(x, y, z).Name == "Snowy_Grass" {
			// check if there is air above it
			if chunk.GetVoxel(x, y, z+1).Name != "Air" {
				return false
			}

			// leaves motherfucka
			leaves := "Leaves"
			if chunk.GetVoxel(x, y, z).Name == "Snowy_Grass" {
				leaves = "Snowy_Leaves"
			}
			for x2 := x - 1; x2 < x+2; x2++ {
				for y2 := y - 1; y2 < y+2; y2++ {
					chunk.SetVoxel(x2, y2, z+4, defaultVoxelDictionary.GetVoxelPointerTo(leaves))
					if math.Abs(float64(x2-x)) != math.Abs(float64(y2-y)) {
						chunk.SetVoxel(x2, y2, z+5, defaultVoxelDictionary.GetVoxelPointerTo(leaves))
					}
				}
			}
			chunk.SetVoxel(x, y, z+5, defaultVoxelDictionary.GetVoxelPointerTo(leaves)) // top middle leaf

			// place the trunk
			for z2 := z + 1; z2 <= z+4; z2++ {
				chunk.SetVoxel(x, y, z2, defaultVoxelDictionary.GetVoxelPointerTo("Wood"))
			}
			return true
		}
	}
	return false
}
