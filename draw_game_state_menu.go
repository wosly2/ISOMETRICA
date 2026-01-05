package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func gameStateDrawMenu(game *Game, screen *ebiten.Image) error {
	game.Framebuffer.Fill(color.RGBA{0, 0, 88, 255})

	game.drawString(game.Framebuffer, "ISOMETRICA", 100, 100, true)
	game.drawString(game.Framebuffer, "GAME IS PAUSED", 100, 115, true)
	game.drawString(game.Framebuffer, "Press Enter", 100, 130, true)

	screen.DrawImage(game.Framebuffer, nil)

	return nil
}
