package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (game *Game) Draw(screen *ebiten.Image) {
	// initiate
	if !game.HasInitiatedDraw {
		game.HasInitiatedDraw = true

		// init camera
		game.Camera[0] = float32(screen.Bounds().Dx() / 2)
		game.Camera[1] = float32(screen.Bounds().Dy() / 2)
		game.ScreenX = screen.Bounds().Dx()
		game.ScreenY = screen.Bounds().Dy()

		// create framebuffer
		game.Framebuffer = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	}

	// ensure framebuffer is the same size as the screen
	if game.Framebuffer.Bounds().Dx() != screen.Bounds().Dx() || game.Framebuffer.Bounds().Dy() != screen.Bounds().Dy() {
		game.Framebuffer = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	}

	// game state
	switch game.GameState {
	case GAMESTATE_GAME:
		_ = gameStateDrawRun(game, screen)
	case GAMESTATE_TITLE:
		_ = gameStateDrawTitle(game, screen)
	case GAMESTATE_MENU:
		_ = gameStateDrawMenu(game, screen)
	}

}
