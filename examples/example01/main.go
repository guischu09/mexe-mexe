package main

import (
	"fmt"
	"mexemexe/pkg/engine"
)

func main() {
	gameConfig := engine.GameConfig{
		Seed:        engine.UNIQUE_SHUFFLE_SEED,
		PlayersName: []string{"Guilherme", "Michele"},
		NumPlayers:  2,
		NumCards:    11,
	}

	game := engine.NewGame(gameConfig)
	game.Players[0].Print()
	game.Players[0].PrintHand()

	game.Players[1].Print()
	game.Players[1].PrintHand()
	fmt.Println("=================")

	// Example 1: Valid Sequence Meld
	twoOfDiamonds := engine.TWO_OF_DIAMONDS
	threeOfDiamonds := engine.THREE_OF_DIAMONDS
	fourOfDiamonds := engine.FOUR_OF_DIAMONDS

	sequenceMeldCards := []engine.Card{twoOfDiamonds, threeOfDiamonds, fourOfDiamonds}
	fmt.Println("Testing sequence meld:")
	meld, err := engine.MakeMeldFromCards(sequenceMeldCards)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Sequence meld is valid!")
		meld.Print()
	}
	fmt.Println()

	// Example 2: Valid Book Meld
	twoOfClubs, _ := engine.NewCard(engine.TWO, engine.CLUB, engine.TWO_VALUE, engine.TWO_CLUB_SYMBOL, engine.BLACK)
	twoOfHearts, _ := engine.NewCard(engine.TWO, engine.HEART, engine.TWO_VALUE, engine.TWO_HEART_SYMBOL, engine.RED)

	bookMeldCards := []engine.Card{twoOfDiamonds, twoOfClubs, twoOfHearts}
	fmt.Println("Testing book meld:")
	meld, err = engine.MakeMeldFromCards(bookMeldCards)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Book meld is valid!")
		meld.Print()
	}
	fmt.Println()

	// Example 3: Invalid Meld - Too few cards
	fmt.Println("Testing invalid meld (too few cards):")
	invalidMeldCards := []engine.Card{twoOfDiamonds, threeOfDiamonds}
	meld, err = engine.MakeMeldFromCards(invalidMeldCards)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Meld is somehow valid (unexpected)!")
		meld.Print()
	}
	fmt.Println()

	// Example 4: Invalid Meld - Mixed cards that don't form sequence or book
	sevenOfSpades, _ := engine.NewCard(engine.SEVEN, engine.SPADE, engine.SEVEN_VALUE, engine.SEVEN_SPADE_SYMBOL, engine.BLACK)
	kingOfClubs, _ := engine.NewCard(engine.KING, engine.CLUB, engine.KING_VALUE, engine.KING_CLUB_SYMBOL, engine.BLACK)

	mixedCards := []engine.Card{twoOfDiamonds, sevenOfSpades, kingOfClubs}
	fmt.Println("Testing invalid meld (mixed cards):")
	meld, err = engine.MakeMeldFromCards(mixedCards)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Meld is somehow valid (unexpected)!")
		meld.Print()
	}
	fmt.Println()

	// Example 5: Test player hand for valid melds
	fmt.Println("Testing player's hand for potential melds:")
	playerCards := game.Players[0].Hand.Cards

	// Try different card combinations from player's hand
	// This assumes you have enough cards in hand to potentially form melds
	if len(playerCards) >= 3 {
		// Try first three cards
		firstThreeCards := playerCards[:3]
		fmt.Println("Testing first three cards:")
		engine.PrintCards(firstThreeCards)

		meld, err = engine.MakeMeldFromCards(firstThreeCards)
		if err != nil {
			fmt.Println("Not a valid meld:", err)
		} else {
			fmt.Println("Found valid meld!")
			meld.Print()
		}

		// Sort cards by value and try again with potentially better grouping
		sortedCards := make([]engine.Card, len(playerCards))
		copy(sortedCards, playerCards)
		engine.SortCardsByValue(sortedCards)

		fmt.Println("\nTesting after sorting by value:")
		for i := 0; i <= len(sortedCards)-3; i++ {
			threeCards := sortedCards[i : i+3]
			engine.PrintCards(threeCards)

			meld, err = engine.MakeMeldFromCards(threeCards)
			if err == nil {
				fmt.Println("Found valid meld!")
				meld.Print()
				break
			}
		}
	}
}
