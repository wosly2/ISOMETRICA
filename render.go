package main

import (
	"image"
	"image/color"
	"log"

	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

const tileWidth int = 32
const tileHeight int = 32

var face = font.Face(basicfont.Face7x13)

// init the rendering system
func initRender() {
	var err error
	groundTextureAtlas, _, err = ebitenutil.NewImageFromFile("block_atlas.png")
	if err != nil {
		log.Fatal(err)
	}
}

// draw a string at a given position
func drawString(screen *ebiten.Image, output string, x, y int) {
	text.Draw(screen, output, face, x, y, color.White)
}

// get the screen position of a voxel
func getScreenPosition(x, y, z int, cameraX, cameraY int, depthShake float32) (screenX, screenY int) {
	screenX = ((x - y) * tileWidth / 2) + cameraX
	// depth shake
	if depthShake != 0 {
		screenX += int(float32(depthShake) * float32(x+y) / 2)
	}
	screenY = ((x+y)*tileHeight/4 - z*tileHeight/2) + cameraY
	return
}

// render a chunk with a given camera position
func (chunk Chunk) Render(screen *ebiten.Image, cameraX, cameraY int, depthShake float32) (blocksRendered int) {
	blocksRendered = 0
	// iterate through voxels
	for x := 0; x < chunk.Width; x++ {
		for y := 0; y < chunk.Height; y++ {
			for z := 0; z < chunk.Depth; z++ {

				// // get the screen position
				screenX, screenY := getScreenPosition(x, y, z, cameraX, cameraY, depthShake)

				// don't bother drawing it if it's off screen
				if screenX+tileWidth < 0 || screenX > screen.Bounds().Dx() || screenY+tileHeight < 0 || screenY > screen.Bounds().Dy() {
					continue
				}

				// get the Voxel type
				var currentVoxel = chunk.GetVoxel(x, y, z)
				var voxelDict = chunk.GetVoxelDictionary(x, y, z)
				if currentVoxel.Name == "Air" {
					continue
				}

				// check if the voxel is even visible
				if !slices.Contains(voxelDict.GetTranparentNames(), chunk.GetVoxel(x+1, y, z).Name) &&
					!slices.Contains(voxelDict.GetTranparentNames(), chunk.GetVoxel(x, y+1, z).Name) &&
					!slices.Contains(voxelDict.GetTranparentNames(), chunk.GetVoxel(x, y, z+1).Name) &&
					// let it render if it's on the edge of the chunk
					!(x == 0 || x == chunk.Width-1 || y == 0 || y == chunk.Height-1 || z == 0 || z == chunk.Depth-1) {
					continue // Skip rendering this voxel
				}

				// hide any transparent under itself (only Transparent, not TransparentNoCull)
				if slices.Contains(voxelDict.Transparent, currentVoxel.Name) {
					if chunk.GetVoxel(x+1, y, z).Name == currentVoxel.Name &&
						chunk.GetVoxel(x, y+1, z).Name == currentVoxel.Name &&
						chunk.GetVoxel(x, y, z+1).Name == currentVoxel.Name {
						//!(x == 0 || x == chunk.Width-1 || y == 0 || y == chunk.Height-1 || z == 0 || z == chunk.Depth-1) {
						continue
					}
				}

				// get the texture from the atlas
				if currentVoxel.Atlas == nil {
					log.Fatal("Atlas is nil")
				}
				subImage := currentVoxel.Atlas.SubImage(image.Rect(currentVoxel.Texture[0], currentVoxel.Texture[1], currentVoxel.Texture[2], currentVoxel.Texture[3]))
				texture, ok := subImage.(*ebiten.Image)
				if !ok {
					continue
				}

				// draw the texture
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(screenX), float64(screenY))
				screen.DrawImage(texture, op)

				blocksRendered++
			}
		}
	}
	return
}
