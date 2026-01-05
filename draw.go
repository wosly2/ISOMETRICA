package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (game *Game) Draw(screen *ebiten.Image) {
	// initiate
	if !game.HasInitiatedDraw {
		game.HasInitiatedDraw = true

		// init camera
		game.Camera = [2]float32{0, 0} // will get updated to correct values later
		game.ScreenX = screen.Bounds().Dx()
		game.ScreenY = screen.Bounds().Dy()

		// create framebuffer
		game.Framebuffer = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

		// assign direction
		game.Direction = SOUTH

		// load font
		game.Font = FontLibrary["Isometrica"].newFont()
	}

	// ensure framebuffer is the same size as the screen
	if game.Framebuffer.Bounds().Dx() != screen.Bounds().Dx() || game.Framebuffer.Bounds().Dy() != screen.Bounds().Dy() {
		game.Framebuffer = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
		game.ScreenX = screen.Bounds().Dx()
		game.ScreenY = screen.Bounds().Dy()
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
