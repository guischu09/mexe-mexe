package engine

import (
	"fmt"
	"slices"
)

type Hand struct {
	Cards []*Card
	Size  int
}

func (h *Hand) RemoveCard(card *Card) bool {
	for i := range h.Cards {
		if h.Cards[i].UUID == card.UUID {
			h.Cards = slices.Delete(h.Cards, i, i+1)
			h.updateSize()
			return true
		}
	}
	return false
}

func (h *Hand) Print() {
	printHand := ""
	for i := range h.Cards {
		printHand += string(h.Cards[i].Symbol) + " "
	}
	fmt.Println(printHand)
}

func (h *Hand) AddCard(card *Card) {
	for i := range h.Cards {
		if h.Cards[i].UUID == card.UUID {
			return
		}
	}
	h.Cards = append(h.Cards, card)
	h.updateSize()
}

func (h *Hand) PrintSize() {
	fmt.Println("Size: ", h.Size)
}

func (h *Hand) updateSize() {
	h.Size = len(h.Cards)
}

func (h *Hand) Contains(card *Card) bool {
	for i := range h.Cards {
		if h.Cards[i].UUID == card.UUID {
			return true
		}
	}
	return false
}

func NewHandFromDeck(deck *Deck, numCards uint8) *Hand {
	cards := make([]*Card, numCards)
	for i := range int(numCards) {
		cards[i] = deck.DrawCard()
	}
	hand := Hand{
		Cards: cards,
		Size:  len(cards),
	}
	return &hand
}

func NewHandFromCards(cards []*Card) *Hand {
	hand := Hand{
		Cards: cards,
		Size:  len(cards),
	}
	return &hand
}
