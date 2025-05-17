package engine

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"slices"

	"golang.org/x/term"
)

type Rederer struct {
	Width         int
	Table         *Table
	Hand          *Hand
	PlayerName    string
	currentPos    int
	selectedCount int
	selectedCards []bool
	turnState     *TurnState
}

func NewRenderer(table *Table, hand *Hand, playerName string, turnState *TurnState) *Rederer {
	return &Rederer{
		Width:         0,
		Table:         table,
		Hand:          hand,
		PlayerName:    playerName,
		currentPos:    0,
		selectedCount: 0,
		selectedCards: []bool{false},
		turnState:     turnState,
	}
}

func (r *Rederer) CreateHorizontalLine(char string) string {
	line := ""
	for i := 0; i < r.Width; i++ {
		line += char
	}
	return line
}

func (r *Rederer) PrintInstructions(screenBuffer *strings.Builder) {

	titleText := "INSTRUCTIONS"
	padding := max((r.Width-len(titleText))/2, 0)
	headerLine := r.CreateHorizontalLine("-")
	padStr := ""
	for range padding {
		padStr += "-"
	}
	screenBuffer.WriteString(fmt.Sprintf("%s%s%s\r\n", padStr, titleText, padStr))

	// Instructions line
	instText := "'s': Select | 'p': Play meld | 'q': Quit | 'd': Draw card | 'e': End turn"
	screenBuffer.WriteString(fmt.Sprintf("%s\r\n", instText))
	screenBuffer.WriteString(fmt.Sprintf("%s\r\n", headerLine[:r.Width]))
}

func (r *Rederer) UserInputDisplay() Play {

	fmt.Print("\033[H\033[2J")

	// Get terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Error setting up terminal: %s\r\n", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Get terminal dimensions
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Error getting terminal size: %s\r\n", err)
		return nil
	}
	r.Width = width

	// Combine cards from hand and table for navigation
	allCards := slices.Clone(r.Hand.Cards)
	tableCards := slices.Clone(r.Table.Cards)
	allCards = append(allCards, tableCards...)

	if len(allCards) == 0 {
		log.Fatal("No cards available to select.")
	}

	r.currentPos = 0
	r.selectedCards = make([]bool, len(allCards))
	r.selectedCount = 0
	statusMessage := ""

	for {
		r.RenderScreen(allCards, statusMessage)
		if statusMessage != "" {
			time.Sleep(1 * time.Second)
			statusMessage = ""
			continue
		}

		buffer := make([]byte, 3)
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading input: %s\r\n", err)
			continue
		}

		if n == 1 {
			switch buffer[0] {
			case 'q':
				fmt.Print("\033[H\033[2J")
				fmt.Printf("Are you sure you want to quit? (y/n)\r\n")
				_, err = os.Stdin.Read(buffer[:1])
				if err != nil {
					continue
				}
				if strings.ToLower(string(buffer[0])) == "y" || strings.ToLower(string(buffer[0])) == "q" {
					return NewQuitPlay("q")
				} else {
					fmt.Print("\033[H\033[2J")
					continue
				}
			case 's':
				r.selectedCards[r.currentPos] = !r.selectedCards[r.currentPos]
				if r.selectedCards[r.currentPos] {
					r.selectedCount++
				} else {
					r.selectedCount--
				}

			case 'd':
				if !r.turnState.HasDrawedCard {
					return NewDrawCardPlay("d")
				} else {
					statusMessage = "You can't draw a card twice in a turn."
					continue
				}

			case 'e':
				if r.turnState.HasPlayedMeld || r.turnState.HasDrawedCard {
					return NewEndTurnPlay("e")
				} else {
					statusMessage = "You must play a meld or draw a card before ending the turn."
					continue
				}

			case 'p':
				if r.selectedCount < MIN_MELD_SIZE {
					statusMessage = fmt.Sprintf("ERROR: You need to select at least %d cards for a meld.", MIN_MELD_SIZE)
					continue
				}

				var selectedMeldCards []*Card
				var handCardUUIDs = make(map[uint8]bool) // Track hand card UUIDs

				// First, build a map of all hand card UUIDs for quick lookup
				for _, card := range r.Hand.Cards {
					handCardUUIDs[card.UUID] = true
				}

				// Track which selected cards are from hand vs table
				var handSelectedCards []*Card
				var tableSelectedCards []*Card

				for i, isSelected := range r.selectedCards {
					if isSelected {
						card := allCards[i]
						selectedMeldCards = append(selectedMeldCards, card)

						// Check if this card is from the hand by UUID
						if handCardUUIDs[card.UUID] {
							handSelectedCards = append(handSelectedCards, card)
						} else {
							tableSelectedCards = append(tableSelectedCards, card)
						}
					}
				}

				// Validate the combined meld
				_, err := MakeMeldFromCards(selectedMeldCards)
				if err != nil {
					statusMessage = fmt.Sprintf("%s", err)
					continue
				}

				fmt.Print("\033[?25h")     // Show cursor before returning
				fmt.Print("\033[H\033[2J") // Clear screen

				// Only pass the hand cards to be moved
				return NewMeldPlay("m", handSelectedCards)
			}
		} else if n == 3 && buffer[0] == 27 && buffer[1] == 91 {
			switch buffer[2] {
			case 68: // Left arrow
				if r.currentPos > 0 {
					r.currentPos--
				}
			case 67: // Right arrow
				if r.currentPos < len(allCards)-1 {
					r.currentPos++
				}
			}
		}
	}
}

func (r *Rederer) RenderScreen(allCards []*Card, statusMessage string) {

	var screenBuffer strings.Builder

	// Display table section
	tableTitle := "\nTABLE"
	tablePadding := (r.Width - len(tableTitle)) / 2
	if tablePadding < 0 {
		tablePadding = 0
	}

	tablePadStr := ""
	for range tablePadding {
		tablePadStr += " "
	}

	screenBuffer.WriteString(fmt.Sprintf("\r\n%s%s\r\n", tablePadStr, tableTitle))
	screenBuffer.WriteString(fmt.Sprintf("%s\r\n", r.CreateHorizontalLine("_")[:r.Width]))

	if len(r.Table.Cards) > 0 {
		tableOffset := len(r.Hand.Cards)
		tableOutput := DisplayCardsWithSelectionToString(r.Table.Cards, r.selectedCards[tableOffset:],
			r.currentPos >= tableOffset, r.currentPos-tableOffset)
		screenBuffer.WriteString(tableOutput)
		screenBuffer.WriteString("\r\n\r\n")
	} else {
		screenBuffer.WriteString("\r\n\r\n")
	}

	screenBuffer.WriteString(fmt.Sprintf("%s\r\n", r.CreateHorizontalLine("_")[:r.Width]))

	// Display hand section
	handTitle := fmt.Sprintf("%s's hand", r.PlayerName)

	handPadding := max((r.Width-len(handTitle))/2, 0)

	handPadStr := ""
	for range handPadding {
		handPadStr += " "
	}

	screenBuffer.WriteString(fmt.Sprintf("\r\n%s%s\r\n", handPadStr, handTitle))

	handDivider := r.CreateHorizontalLine("-")[:r.Width]
	screenBuffer.WriteString(fmt.Sprintf("\r\n%s\r\n", handDivider))

	handOutput := DisplayCardsWithSelectionToString(r.Hand.Cards, r.selectedCards[:len(r.Hand.Cards)],
		r.currentPos < len(r.Hand.Cards), r.currentPos)
	screenBuffer.WriteString(handOutput)

	screenBuffer.WriteString(fmt.Sprintf("\r\n%s\r\n", handDivider))

	// Show selected cards
	if r.selectedCount > 0 {
		selectedCardsList := []Card{}
		for i, isSelected := range r.selectedCards {
			if isSelected {
				selectedCardsList = append(selectedCardsList, *allCards[i])
			}
		}
		selectedOutput := DisplayCardsToString(selectedCardsList, -1)
		screenBuffer.WriteString(selectedOutput)
	}

	if statusMessage != "" {
		screenBuffer.WriteString(fmt.Sprintf("\r\n\r\n%s\r\n", statusMessage))
	}

	fmt.Print("\033[H")
	fmt.Print(screenBuffer.String())
	fmt.Print("\033[J")

}

// New helper functions that return strings instead of printing directly
func DisplayCardsWithSelectionToString(cards []*Card, selections []bool, isCurrentSection bool, currentPosInSection int) string {
	var output strings.Builder

	if len(cards) == 0 {
		output.WriteString("No cards.\r\n")
		return output.String()
	}

	output.WriteString("  ")
	for i, card := range cards {
		if isCurrentSection && i == currentPosInSection {
			output.WriteString("\033[7m")
		}
		if selections[i] {
			output.WriteString("*")
		}

		// Print card
		output.WriteString(fmt.Sprintf("%s ", string(card.Symbol)))

		if isCurrentSection && i == currentPosInSection {
			output.WriteString("\033[0m")
		}
	}

	return output.String()
}

// Helper function to display cards without selection (for showing selected cards)
func DisplayCardsToString(cards []Card, highlightPos int) string {
	var output strings.Builder

	if len(cards) == 0 {
		output.WriteString("No cards.\r\n")
		return output.String()
	}

	for i, card := range cards {
		if i == highlightPos {
			output.WriteString("\033[7m")
		}

		output.WriteString(fmt.Sprintf("%s ", string(card.Symbol)))

		if i == highlightPos {
			output.WriteString("\033[0m")
		}
	}

	return output.String()
}
