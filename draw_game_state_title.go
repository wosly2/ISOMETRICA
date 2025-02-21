package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func gameStateDrawTitle(game *Game, screen *ebiten.Image) error {
	game.Framebuffer.Fill(color.RGBA{0, 0, 88, 255})

	drawString(game.Framebuffer, "ISOMETRICA", 100, 100, true)
	drawString(game.Framebuffer, "COPYRIGHT MMXXV SOYPACKET", 100, 115, true)
	drawString(game.Framebuffer, "Press Any Key!", 100, 145, true)
	drawString(game.Framebuffer, fmt.Sprintf("%d", game.World.Seed), 100, 160, true)

	screen.DrawImage(game.Framebuffer, nil)

	return nil
}
