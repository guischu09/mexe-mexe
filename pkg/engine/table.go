package engine

import (
	"fmt"
	"slices"
)

type Table struct {
	Cards []*Card
	Size  int
}

func (t *Table) RemoveCard(card *Card) bool {
	for i := 0; i < len(t.Cards); i++ {
		if t.Cards[i].Name == card.Name {
			t.Cards = slices.Delete(t.Cards, i, i+1)
			t.updateSize()
			return true
		}
	}
	fmt.Println("Card not found in the Table")
	return false
}

func (t *Table) Print() {
	printTable := ""
	for i := 0; i < len(t.Cards); i++ {
		printTable += string(t.Cards[i].Symbol) + " "
	}
	fmt.Println(printTable)
}

func (t *Table) AddCard(card *Card) {
	t.Cards = append(t.Cards, card)
	t.updateSize()
}

func (t *Table) updateSize() {
	t.Size = len(t.Cards)
}
