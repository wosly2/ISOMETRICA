package main

import (
	"image/color"
	"log"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileWidth  int = 32
	tileHeight int = 32
	v          int = 32
)

// fuck non constant slices
// don't assign to these
var (
	// sxX, sxY, syX, syY
	SOUTH = [4]int{1, 1, -1, 1}
	WEST  = [4]int{-1, 1, -1, -1}
	NORTH = [4]int{-1, -1, 1, -1}
	EAST  = [4]int{1, -1, 1, 1}
)

// init the rendering system
func initRender() {
	var err error
	groundTextureAtlas, _, err = ebitenutil.NewImageFromFile("assets/block_atlas.png")
	if err != nil {
		log.Fatal(err)
	}
}

// draw a string at a given position
func (game *Game) drawString(screen *ebiten.Image, output string, x, y int, shadow bool) {
	if shadow {
		game.Font.renderString(screen, output, x+1, y+1, color.RGBA{0, 0, 0, 255})
	}
	game.Font.renderString(screen, output, x, y, color.RGBA{255, 255, 255, 255})
}

// get the screen position of a voxel
func getScreenPosition(x, y, z int, cameraX, cameraY float32, depthShake float32, direction [4]int) (screenX, screenY int) {
	sxX, sxY, syX, syY := direction[0], direction[1], direction[2], direction[3]

	screenX = ((x*sxX + y*syX) * v / 2)
	screenY = ((x*sxY + y*syY) * v / 4) - z*v/2

	screenX += int(cameraX)
	screenY += int(cameraY)

	if depthShake != 0 {
		localX := x % 64
		localY := y % 64
		screenX += int(float32(depthShake) * float32(math.Sin(float64(localX+localY))))
	}

	return
}

func (chunk *Chunk) VoxelIsVisible(x, y, z int) bool {
	// check if voxel is in bounds
	if x < 0 || y < 0 || z < 0 || x >= chunk.Width || y >= chunk.Height || z >= chunk.Depth {
		return false
	}

	var voxelDict = chunk.GetVoxelDictionary(x, y, z)
	return slices.Contains(voxelDict.GetTransparentNames(), chunk.GetVoxel(x+1, y, z).Name) ||
		slices.Contains(voxelDict.GetTransparentNames(), chunk.GetVoxel(x, y+1, z).Name) ||
		slices.Contains(voxelDict.GetTransparentNames(), chunk.GetVoxel(x, y, z+1).Name) ||
		// let it render if it's on the edge of the chunk
		(x == 0 || x == chunk.Width-1 || y == 0 || y == chunk.Height-1 || z == 0 || z == chunk.Depth-1)
}

func ChunkContainingGlobalPointVisibleInViewport(x, y, z int, cameraX, cameraY float32, depthShake float32, screenWidth, screenHeight int, direction [4]int) bool {
	// get the voxel space origin of the chunk
	var chunkX, chunkY, chunkZ int = (x / 32) * 32, (y / 32) * 32, (z / 32) * 32

	// get the screen space bounds (diamond vertices) of the chunk
	diamondTopX, diamondTopY := getScreenPosition(chunkX, chunkY, chunkZ+32, cameraX, cameraY, depthShake, direction)
	diamondBottomX, diamondBottomY := getScreenPosition(chunkX+32, chunkY+32, chunkZ, cameraX, cameraY, depthShake, direction)
	diamondLeftX, diamondLeftY := getScreenPosition(chunkX, chunkY+32, chunkZ+32, cameraX, cameraY, depthShake, direction)
	diamondRightX, diamondRightY := getScreenPosition(chunkX+32, chunkY, chunkZ+32, cameraX, cameraY, depthShake, direction)

	// check if any lines of the diamond intersect the camera
	// top to right
	if doesLineIntersectRectangle(diamondTopX, diamondTopY, diamondRightX, diamondRightY, 0, 0, screenWidth, screenHeight) {
		return true
	}
	// right to bottom
	if doesLineIntersectRectangle(diamondRightX, diamondRightY, diamondBottomX, diamondBottomY, 0, 0, screenWidth, screenHeight) {
		return true
	}
	// bottom to left
	if doesLineIntersectRectangle(diamondBottomX, diamondBottomY, diamondLeftX, diamondLeftY, 0, 0, screenWidth, screenHeight) {
		return true
	}
	// left to top
	if doesLineIntersectRectangle(diamondLeftX, diamondLeftY, diamondTopX, diamondTopY, 0, 0, screenWidth, screenHeight) {
		return true
	}

	// must not be in bounds
	return false
}

// render a chunk with a given camera position
func (chunk Chunk) Render(screen *ebiten.Image, cameraX, cameraY float32, depthShake float32, game *Game, renderPlayer bool) (blocksRendered int) {
	// create an offscreen buffer to render to
	chunkRenderBuffer := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	// cache some values for performance
	chunkWidth := chunk.Width
	chunkHeight := chunk.Height
	chunkDepth := chunk.Depth
	voxelDict := defaultVoxelDictionary
	transparentNames := voxelDict.GetTransparentNames()
	screenWidth := screen.Bounds().Dx()
	screenHeight := screen.Bounds().Dy()

	// get the proper render order based on camera direction
	startX, stopX, stepX := 0, 0, 0
	startY, stopY, stepY := 0, 0, 0
	startZ, stopZ, stepZ := 0, 0, 0
	switch game.Direction {
	case SOUTH: // CORRECT
		startX, stopX, stepX = 0, chunkWidth, 1
		startY, stopY, stepY = 0, chunkHeight, 1
		startZ, stopZ, stepZ = 0, chunkDepth, 1
	case NORTH: // FIXME
		startX, stopX, stepX = chunkWidth-1, 0, -1
		startY, stopY, stepY = chunkHeight-1, 0, -1
		startZ, stopZ, stepZ = 0, chunkDepth, 1
	case EAST: // FIXME
		startX, stopX, stepX = 0, chunkWidth, 1
		startY, stopY, stepY = chunkHeight-1, 0, -1
		startZ, stopZ, stepZ = 0, chunkDepth, 1
	case WEST: // CORRECT
		startX, stopX, stepX = 0, chunkWidth, 1
		startY, stopY, stepY = chunkHeight-1, 0, -1
		startZ, stopZ, stepZ = 0, chunkDepth, 1
	}

	blocksRendered = 0
	// iterate through voxels
	for x := startX; x != stopX; x += stepX {
		for y := startY; y != stopY; y += stepY {
			for z := startZ; z != stopZ; z += stepZ {
				// // get the screen position
				screenX, screenY := getScreenPosition(x, y, z, cameraX, cameraY, depthShake, game.Direction)

				// don't bother drawing it if it's off screen
				if screenX+tileWidth < 0 || screenX > screenWidth || screenY+tileHeight < 0 || screenY > screenHeight {
					continue
				}

				// get the Voxel type
				var currentVoxel = chunk.GetVoxel(x, y, z)
				if currentVoxel.Name == "Air" {
					continue
				}

				// check if the voxel is even visible
				if !chunk.VoxelIsVisible(x, y, z) {
					continue // Skip rendering this voxel
				}

				// hide any transparent under itself (only Transparent, not TransparentNoCull)
				if slices.Contains(transparentNames, currentVoxel.Name) {
					if chunk.GetVoxel(x+1, y, z).Name == currentVoxel.Name &&
						chunk.GetVoxel(x, y+1, z).Name == currentVoxel.Name &&
						chunk.GetVoxel(x, y, z+1).Name == currentVoxel.Name {
						//!(x == 0 || x == chunk.Width-1 || y == 0 || y == chunk.Height-1 || z == 0 || z == chunk.Depth-1) { // old edge culling, no longer used
						continue
					}
				}

				// get the texture from the atlas
				if currentVoxel.Atlas == nil {
					log.Fatal("Atlas is nil")
				}
				subImage := currentVoxel.Atlas.SubImage(currentVoxel.TextureRect)
				texture, ok := subImage.(*ebiten.Image)
				if !ok {
					continue
				}

				// draw the texture
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(screenX), float64(screenY))
				chunkRenderBuffer.DrawImage(texture, op)

				blocksRendered++
			}
		}
	}

	// render the offscreen buffer to the screen
	screen.DrawImage(chunkRenderBuffer, nil)

	return
}
