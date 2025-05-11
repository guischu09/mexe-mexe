package main

import "mexemexe/internal/engine"

//TODO: Implement a mechanic to control validation of played melds and melds remaining on the table.
// TODO: Implement validation of Books

func main() {

	gameConfig := engine.GameConfig{
		Seed:              engine.UNIQUE_SHUFFLE_SEED,
		PlayersName:       []string{"Guilherme", "Michele"},
		NumPlayers:        2,
		NumCards:          21,
		RandomPlayerOrder: true,
		TotalCards:        uint8(engine.TOTAL_DECK_SIZE),
	}

	game := engine.NewGame(gameConfig)
	if !game.Start() {
		game.Close()
	}

}
