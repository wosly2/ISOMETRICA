package main

import "github.com/hajimehoshi/ebiten/v2"

func gameStateUpdateTitle(game *Game) error {
	if len(ebiten.InputChars()) > 0 {
		game.GameState = GAMESTATE_GAME
	}

	return nil
}
