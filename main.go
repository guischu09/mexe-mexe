package main

import (
	"fmt"
	"mexemexe/pkg/engine"
)

type GameConfig struct {
	Seed        uint64
	PlayersName []string
	NumPlayers  uint8
	NumCards    uint8
}

type Game struct {
	Config  GameConfig
	Deck    engine.GameDeck
	Table   engine.Table
	Players []engine.Player
}

func NewGame(config GameConfig) Game {

	players := make([]engine.Player, config.NumPlayers)
	deck := engine.NewGameDeck(config.Seed)

	for i := 0; i < int(config.NumPlayers); i++ {
		newHand := engine.NewHandFromDeck(&deck, config.NumCards)
		players[i] = engine.NewPlayer(config.PlayersName[i], newHand, 0)
	}

	return Game{
		Config:  config,
		Deck:    deck,
		Players: players,
	}
}

func main() {

	gameConfig := GameConfig{
		Seed:        engine.UNIQUE_SHUFFLE_SEED,
		PlayersName: []string{"Guilherme", "Michele"},
		NumPlayers:  2,
		NumCards:    11,
	}

	game := NewGame(gameConfig)
	game.Players[0].Print()
	game.Players[0].PrintHand()

	game.Players[1].Print()
	game.Players[1].PrintHand()
	fmt.Println("=================")
	// cards := game.Players[0].Hand.Cards
	// fmt.Println("Cards:")
	// engine.PrintCards(cards)

	// engine.SortCardsByValue(cards)
	// fmt.Println("Sorted cards by value:")
	// engine.PrintCards(cards)

	// engine.SortCardsBySuitAndValue(cards)
	// fmt.Println("Sorted cards by suit and value:")
	// engine.PrintCards(cards)

	twoOfDiamonds, _ := engine.NewCard(engine.TWO, engine.DIAMOND, engine.TWO_VALUE, engine.TWO_DIAMOND_SYMBOL, engine.RED)
	threeOfDiamonds, _ := engine.NewCard(engine.THREE, engine.DIAMOND, engine.THREE_VALUE, engine.THREE_DIAMOND_SYMBOL, engine.RED)
	fourOfDiamonds, _ := engine.NewCard(engine.FOUR, engine.DIAMOND, engine.FOUR_VALUE, engine.FOUR_DIAMOND_SYMBOL, engine.RED)

	attemptMeldCards := []engine.Card{twoOfDiamonds, threeOfDiamonds, fourOfDiamonds}

	meld, err := engine.MakeMeldFromCards(attemptMeldCards)
	if err != nil {
		fmt.Println("ERROR: Not a valid Meld")
	} else {
		fmt.Println("Meld is valid!")
		meld.Print()
	}

	attemptMeldCards2 := []engine.Card{}

	meld, err := engine.MakeMeldFromCards(attemptMeldCards)
	if err != nil {
		fmt.Println("ERROR: Not a valid Meld")
	} else {
		fmt.Println("Meld is valid!")
		meld.Print()
	}

	// if game.Deck.Contains(game.Players[0].Hand.Cards[0]) {
	// fmt.Printf("Card %s found in deck\n", game.Players[0].Hand.Cards[0].Symbol)
	// } else {
	// fmt.Printf("Card %s not found in deck\n", game.Players[0].Hand.Cards[0].Symbol)
	// }

}
