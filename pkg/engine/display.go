package engine

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

func MeldDisplayInput(table *Table, hand *Hand) []Card {
	// Setup terminal for keyboard input
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error setting up terminal:", err)
		return nil
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Combine cards from hand and table for navigation
	allCards := append([]Card{}, hand.Cards...)
	tableCards := append([]Card{}, table.PlayedCards...)
	allCards = append(allCards, tableCards...)

	if len(allCards) == 0 {
		fmt.Println("No cards available to select.")
		return nil
	}

	// Track current position and selections
	currentPos := 0
	selectedCards := make([]bool, len(allCards))
	selectedCount := 0

	// Main interaction loop
	for {
		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Display instructions
		fmt.Printf("========================== MELD SELECTION MODE ===================\r\n")
		fmt.Printf("  ← → : Navigate | 's': Select | 'p': Play meld | 'q': Quit\r\n")
		fmt.Printf("==================================================================\r\n")

		// Display table section
		fmt.Print("\r\n                                TABLE\r\n")
		fmt.Printf("__________________________________________________________________\r\n")
		if len(tableCards) > 0 {
			tableOffset := len(hand.Cards)
			DisplayCardsWithSelection(tableCards, selectedCards[tableOffset:], currentPos >= tableOffset, currentPos-tableOffset)
		} else {
			fmt.Printf("\r\n\r\n")
		}

		fmt.Printf("__________________________________________________________________\r\n")
		// Display hand section
		fmt.Printf("\r\n------------------------------ YOUR HAND -------------------------\r\n")
		DisplayCardsWithSelection(hand.Cards, selectedCards[:len(hand.Cards)], currentPos < len(hand.Cards), currentPos)
		fmt.Printf("\r\n------------------------------------------------------------------\r\n")
		// Show selected cards
		if selectedCount > 0 {
			selectedCardsList := []Card{}
			for i, isSelected := range selectedCards {
				if isSelected {
					selectedCardsList = append(selectedCardsList, allCards[i])
				}
			}
			DisplayCards(selectedCardsList, -1)
		}

		// Read key input
		buffer := make([]byte, 3)
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading input: %s\r\n", err)
			continue
		}

		if n == 1 {
			switch buffer[0] {
			case 'q':
				fmt.Printf("Exiting meld selection.\r\n")
				return nil
			case 's':
				// Toggle selection of current card
				selectedCards[currentPos] = !selectedCards[currentPos]
				if selectedCards[currentPos] {
					selectedCount++
				} else {
					selectedCount--
				}
			case 'p':
				if selectedCount < MIN_MELD_SIZE {
					fmt.Printf("Not enough cards selected. You need at least %d cards for a meld.\r\n", MIN_MELD_SIZE)
					time.Sleep(1 * time.Second)
					continue
				}

				var selectedMeldCards []Card
				for i, isSelected := range selectedCards {
					if isSelected {
						selectedMeldCards = append(selectedMeldCards, allCards[i])
					}
				}

				meld, err := MakeMeldFromCards(selectedMeldCards)
				if err != nil {
					time.Sleep(1 * time.Second)
					continue
				}

				fmt.Printf("\r\nValid %s meld created!\r\n", meld.Type)
				return selectedMeldCards
			}
		} else if n == 3 && buffer[0] == 27 && buffer[1] == 91 {
			switch buffer[2] {
			case 68: // Left arrow
				if currentPos > 0 {
					currentPos--
				}
			case 67: // Right arrow
				if currentPos < len(allCards)-1 {
					currentPos++
				}
			}
		}
	}
}

// Helper function to display cards with selection and highlight
func DisplayCardsWithSelection(cards []Card, selections []bool, isCurrentSection bool, currentPosInSection int) {
	if len(cards) == 0 {
		fmt.Printf("No cards.\r\n")
		return
	}

	fmt.Print("  ")
	for i, card := range cards {
		if isCurrentSection && i == currentPosInSection {
			fmt.Print("\033[7m")
		}
		if selections[i] {
			fmt.Print("*")
		}

		// Print card
		fmt.Printf("%s ", string(card.Symbol))

		if isCurrentSection && i == currentPosInSection {
			fmt.Print("\033[0m")
		}
	}

}

// Helper function to display cards without selection (for showing selected cards)
func DisplayCards(cards []Card, highlightPos int) {
	if len(cards) == 0 {
		fmt.Printf("No cards.\r\n")
		return
	}

	for i, card := range cards {
		if i == highlightPos {
			fmt.Print("\033[7m")
		}

		fmt.Printf("%s ", string(card.Symbol))

		if i == highlightPos {
			fmt.Print("\033[0m")
		}
	}
}
