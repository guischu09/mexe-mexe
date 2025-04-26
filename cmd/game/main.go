package main

import "mexemexe/pkg/engine"

func main() {

	gameConfig := engine.GameConfig{
		Seed:        engine.UNIQUE_SHUFFLE_SEED,
		PlayersName: []string{"Guilherme", "Michele"},
		NumPlayers:  2,
		NumCards:    11,
	}

	game := engine.NewGame(gameConfig)

	game.Start()
	// game.Players[0].Print()
	// game.Players[0].PrintHand()

	// game.Start()

}
