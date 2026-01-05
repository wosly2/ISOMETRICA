package main

func (game *Game) Update() error {
	// initiate
	if !game.HasInitiatedUpdate {
		game.HasInitiatedUpdate = true

		go game.syncWorldWithDisk()
	}

	var err error

	// game state
	switch game.GameState {
	case GAMESTATE_GAME:
		err = gameStateUpdateRun(game)
	case GAMESTATE_TITLE:
		err = gameStateUpdateTitle(game)
	case GAMESTATE_MENU:
		err = gameStateUpdateMenu(game)
	}

	return err
}
