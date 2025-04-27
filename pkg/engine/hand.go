package engine

import "fmt"

type Hand struct {
	Cards []Card
	Size  int
}

func (h *Hand) RemoveCard(card Card) bool {
	for i := 0; i < len(h.Cards); i++ {
		if h.Cards[i].Name == card.Name {
			h.Cards = append(h.Cards[:i], h.Cards[i+1:]...)
			h.updateSize()
			return true
		}
	}
	fmt.Println("Card not found in the hand")
	return false
}

func (h *Hand) Print() {
	printHand := ""
	for i := 0; i < len(h.Cards); i++ {
		printHand += string(h.Cards[i].Symbol) + " "
	}
	fmt.Println(printHand)
}

func (h *Hand) AddCard(card Card) {
	h.Cards = append(h.Cards, card)
	h.updateSize()
}

func (h *Hand) DrawCard() Card {
	card := h.Cards[0]
	h.Cards = h.Cards[1:]
	h.updateSize()
	return card
}

func (h *Hand) PrintSize() {
	fmt.Println("Size: ", h.Size)
}

func (h *Hand) updateSize() {
	h.Size = len(h.Cards)
}

func (h *Hand) Contains(card Card) bool {
	for i := 0; i < len(h.Cards); i++ {
		if h.Cards[i].Name == card.Name {
			return true
		}
	}
	return false
}

func NewHandFromDeck(deck *GameDeck, numCards uint8) Hand {
	cards := make([]Card, numCards)
	for i := 0; i < int(numCards); i++ {
		cards[i] = deck.DrawCard()
	}
	hand := Hand{
		Cards: cards,
		Size:  len(cards),
	}
	return hand
}

func NewHandFromCards(cards []Card) Hand {
	hand := Hand{
		Cards: cards,
		Size:  len(cards),
	}
	return hand
}
