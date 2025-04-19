package engine

import "fmt"

type MeldType string

const MIN_MELD_SIZE = 3

const SEQUENCE MeldType = "SEQUENCE"
const BOOK MeldType = "BOOK"
const NONE MeldType = "NONE"

type Meld struct {
	Type  MeldType
	Cards []Card
}

func (m *Meld) Print() {
	printHand := ""
	for i := 0; i < len(m.Cards); i++ {
		printHand += string(m.Cards[i].Symbol) + " "
	}
	fmt.Println(printHand)
}

func (m *Meld) IsValid() bool {
	switch m.Type {
	case SEQUENCE:
		return isMeldSequence(m.Cards)
	case BOOK:
		return isMeldBook(m.Cards)
	default:
		return false
	}
}

func MakeMeldFromCards(cards []Card) (Meld, error) {

	// A meld must have at least 3 cards
	if len(cards) < MIN_MELD_SIZE {
		fmt.Println("ERROR: Meld must have at least 3 cards")
		return Meld{
			Type:  NONE,
			Cards: []Card{},
		}, fmt.Errorf("ERROR: Meld must have at least 3 cards")
	}

	// Check if cards are a sequence (e.g. Q,K,A of hearts)
	SortCardsByValue(cards)

	// Check if cards are a book (e.g. Q, Q, Q)
	if isMeldBook(cards) {
		return Meld{BOOK, cards}, nil
	} else if isMeldSequence(cards) {
		return Meld{SEQUENCE, cards}, nil
	} else {
		fmt.Println("ERROR: Not a valid Meld")
		return Meld{
			Type:  NONE,
			Cards: []Card{},
		}, fmt.Errorf("ERROR: Not a valid Meld")

	}

}

func isMeldSequence(cards []Card) bool {

	// Check if has same suit - if not, it is not a sequence
	for i := 0; i < len(cards)-1; i++ {
		if cards[i].Suit != cards[i+1].Suit {
			return false
		}
	}

	// Check if they are in ascending sequence (e.g. 4,5,6)
	for i := 0; i < len(cards)-1; i++ {
		if cards[i].Value != cards[i+1].Value-1 {
			return false
		}
	}

	return true
}

func isMeldBook(cards []Card) bool {
	// Check if they have the same value (e.g. 4,4,4)
	for i := 0; i < len(cards)-1; i++ {
		if cards[i].Value != cards[i+1].Value {
			return false
		}
	}
	return true
}
