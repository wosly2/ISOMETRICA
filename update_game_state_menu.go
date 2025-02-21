package main

import "github.com/hajimehoshi/ebiten/v2"

func gameStateUpdateMenu(game *Game) error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		game.GameState = GAMESTATE_GAME
	}

	return nil
}
