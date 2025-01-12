package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// player texture atlas
var playerTextureAtlas, _, _ = ebitenutil.NewImageFromFile("player_atlas.png")

// player texture map
// player sprites are 32x48 pixels.
var playerTextureMap = map[string][4]int{
	"Default": [4]int{0, 0, 32, 48},
}

// Player, contains information about a player.
type Player struct {
	Position Vec3
	Velocity Vec3
	Drag     Vec3
	Texture  string
}

var Gravity float32 = 0.01

func (player *Player) Update(world World) {
	// drag
	player.Velocity.X *= player.Drag.X
	player.Velocity.Y *= player.Drag.Y
	player.Velocity.Z *= player.Drag.Z

	// move the player
	player.Position.X += player.Velocity.X
	player.Position.Y += player.Velocity.Y
	player.Position.Z += player.Velocity.Z

	// gravity
	//player.Velocity.Z -= Gravity
}

// get screen position of the player regardless of camera position
func (player *Player) getScreenPosition(depthShake float32) (screenX, screenY float32) {
	screenX = ((player.Position.X - player.Position.Y) * float32(tileWidth) / 2) + float32(tileWidth/2)
	screenY = ((player.Position.X+player.Position.Y)*float32(tileHeight)/4 - (player.Position.Z * float32(tileHeight) / 2)) + float32(tileHeight/2)
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
