package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// player texture atlas
var playerTextureAtlas, _, pErr = ebitenutil.NewImageFromFile("player_atlas.png")

// player texture map
// player sprites are 32x48 pixels.
var playerTextureMap = map[string][4]int{
	"Default": [4]int{0, 0, 32, 48},
}

// Player, contains information about a player.
type Player struct {
	Position [3]float32
	Velocity [3]float32
	Drag     [3]float32
	Texture  string
}

var Gravity float32 = 0.01

func (player *Player) Update(world World) {

	// move the player
	player.Position[0] += player.Velocity[0]
	player.Position[1] += player.Velocity[1]
	player.Position[2] += player.Velocity[2]

	// drag
	player.Velocity[0] *= player.Drag[0]
	player.Velocity[1] *= player.Drag[1]
	player.Velocity[2] *= player.Drag[2]

	// gravity
	//player.Velocity[2] -= Gravity

	// check if the player is on the ground (on top of a voxel)
	var playerVoxelPosition = [3]int{int(player.Position[0]), int(player.Position[1]), int(player.Position[2])}
	voxel, _ := world.GetVoxel(playerVoxelPosition[0], playerVoxelPosition[1], playerVoxelPosition[2])
	if voxel.Name != "Air" {
		player.Velocity[2] = 0 // stop falling
		// move the player to the top of the ground voxel
		player.Position[2] = float32(playerVoxelPosition[2])
	}
}

// get screen position of the player
func (player *Player) getScreenPosition(depthShake float32) (screenX, screenY float32) {
	screenX = ((player.Position[0] - player.Position[1]) * float32(tileWidth) / 2)
	// depth shake
	if depthShake != 0 {
		screenX += (float32(depthShake) * float32(player.Position[0]+player.Position[1]) / 2)
	}
	screenY = ((player.Position[0]+player.Position[1])*float32(tileHeight)/4 - (player.Position[2] * float32(tileHeight) / 2))
	return
}

func (player *Player) Render(screen *ebiten.Image, depthShake float32) {
	// get the texture
	var texture = playerTextureAtlas.SubImage(image.Rect(playerTextureMap[player.Texture][0], playerTextureMap[player.Texture][1], playerTextureMap[player.Texture][2], playerTextureMap[player.Texture][3])).(*ebiten.Image)

	// render
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screen.Bounds().Dx()/2), float64(screen.Bounds().Dy()/2))
	screen.DrawImage(texture, op)
}
