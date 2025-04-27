package main

import (
	"fmt"
	"mexemexe/pkg/engine"
)

func mainasd() {
	gameDeck := engine.NewDeck(engine.UNIQUE_SHUFFLE_SEED)
	gameDeck.PrintSize()
	card := gameDeck.DrawCard()
	card.Print()
	gameDeck.PrintSize()

	card2 := gameDeck.DrawCard()
	card2.Print()
	gameDeck.PrintSize()

	gameDeck.RemoveCard(card2)
	gameDeck.PrintSize()

	if !gameDeck.Contains(card2) {
		fmt.Println("Deck does not contain card")
	} else {
		fmt.Println("Deck contains card.")
	}

}
