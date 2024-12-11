package main

import (
	"image"
	"log"

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

// render a chunk with a given camera position
func (chunk *Chunk) Render(screen *ebiten.Image, cameraX, cameraY int) {
	// iterate through voxels
	for x := 0; x < chunk.Width; x++ {
		for y := 0; y < chunk.Height; y++ {
			for z := 0; z < chunk.Depth; z++ {

				// get the Voxel type
				var currentVoxel = defaultVoxelDictionary.Voxels[chunk.GetVoxel(x, y, z).Index]

				// get the texture from the atlas
				if currentVoxel.Atlas == nil {
					log.Fatal("Atlas is nil")
				}
				subImage := currentVoxel.Atlas.SubImage(image.Rect(currentVoxel.Texture[0], currentVoxel.Texture[1], currentVoxel.Texture[2], currentVoxel.Texture[3]))
				texture, ok := subImage.(*ebiten.Image)
				if !ok {
					continue
				}

				// get the screen position
				screenX := ((x - y) * tileWidth / 2) + cameraX
				screenY := ((x+y)*tileHeight/4 - z*tileHeight/2) + cameraY

				log.Println(screenX, screenY)

				// draw the texture
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(screenX), float64(screenY))
				screen.DrawImage(texture, op)
			}
		}
	}
}
