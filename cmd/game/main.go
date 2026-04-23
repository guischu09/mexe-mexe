package main

import (
	"fmt"
	"mexemexe/internal/engine"
	"mexemexe/internal/service"
)

func main() {

	gameConfig := engine.GameConfig{
		Seed:              engine.UNIQUE_SHUFFLE_SEED,
		PlayersName:       []string{"Gui", "Mi"},
		NumPlayers:        2,
		NumCards:          21,
		RandomPlayerOrder: true,
		TotalCards:        uint8(engine.TOTAL_DECK_SIZE),
	}

	logger := service.NewLogger(0, "test")
	game := engine.NewGame(&gameConfig, logger)
	fmt.Println(game)
}
