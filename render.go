package main

import (
	"image/color"
	"log"

	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const tileWidth int = 32
const tileHeight int = 32

// init the rendering system
func initRender() {
	var err error
	groundTextureAtlas, _, err = ebitenutil.NewImageFromFile("block_atlas.png")
	if err != nil {
		log.Fatal(err)
	}
}

// draw a string at a given position
func drawString(screen *ebiten.Image, output string, x, y int, shadow bool) {
	if shadow {
		isometricaFont.renderString(screen, output, x+1, y+1, color.RGBA{0, 0, 0, 255})
	}
	isometricaFont.renderString(screen, output, x, y, color.RGBA{255, 255, 255, 255})
}

// get the screen position of a voxel
func getScreenPosition(x, y, z int, cameraX, cameraY float32, depthShake float32) (screenX, screenY int) {
	screenX = ((x - y) * tileWidth / 2) + int(cameraX)
	// depth shake
	if depthShake != 0 {
		screenX += int(float32(depthShake) * float32(x+y) / 2)
	}
	screenY = ((x+y)*tileHeight/4 - z*tileHeight/2) + int(cameraY)
	return
}

func (chunk *Chunk) VoxelIsVisible(x, y, z int) bool {
	// check if voxel is in bounds
	if x < 0 || y < 0 || z < 0 || x >= chunk.Width || y >= chunk.Height || z >= chunk.Depth {
		return false
	}

	var voxelDict = chunk.GetVoxelDictionary(x, y, z)
	return slices.Contains(voxelDict.GetTranparentNames(), chunk.GetVoxel(x+1, y, z).Name) ||
		slices.Contains(voxelDict.GetTranparentNames(), chunk.GetVoxel(x, y+1, z).Name) ||
		slices.Contains(voxelDict.GetTranparentNames(), chunk.GetVoxel(x, y, z+1).Name) ||
		// let it render if it's on the edge of the chunk
		(x == 0 || x == chunk.Width-1 || y == 0 || y == chunk.Height-1 || z == 0 || z == chunk.Depth-1)
}

func ChunkContainingGlobalPointVisibleInViewport(x, y, z int, cameraX, cameraY float32, depthShake float32, screenWidth, screenHeight int) bool {
	// get the voxel space origin of the chunk
	var chunkX, chunkY, chunkZ int = (x / 32) * 32, (y / 32) * 32, (z / 32) * 32

	// get the screen space bounds (diamond vertices) of the chunk
	diamondTopX, diamondTopY := getScreenPosition(chunkX, chunkY, chunkZ+32, cameraX, cameraY, depthShake)
	diamondBottomX, diamondBottomY := getScreenPosition(chunkX+32, chunkY+32, chunkZ, cameraX, cameraY, depthShake)
	diamondLeftX, diamondLeftY := getScreenPosition(chunkX, chunkY+32, chunkZ+32, cameraX, cameraY, depthShake)
	diamondRightX, diamondRightY := getScreenPosition(chunkX+32, chunkY, chunkZ+32, cameraX, cameraY, depthShake)

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
	transparentNames := voxelDict.GetTranparentNames()
	screenWidth := screen.Bounds().Dx()
	screenHeight := screen.Bounds().Dy()

	blocksRendered = 0
	// iterate through voxels
	for x := 0; x < chunkWidth; x++ {
		for y := 0; y < chunkHeight; y++ {
			for z := 0; z < chunkDepth; z++ {
				// // get the screen position
				screenX, screenY := getScreenPosition(x, y, z, cameraX, cameraY, depthShake)

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
