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
		if t.Cards[i].UUID == card.UUID {
			t.Cards = slices.Delete(t.Cards, i, i+1)
			t.updateSize()
			return true
		}
	}
	return false
}

func (t *Table) Print() {
	printTable := ""
	for i := 0; i < len(t.Cards); i++ {
		printTable += string(t.Cards[i].Symbol) + " "
	}
	fmt.Println(printTable)
}

func (g *Table) Contains(card *Card) bool {
	for i := 0; i < len(g.Cards); i++ {
		if g.Cards[i].UUID == card.UUID {
			return true
		}
	}
	return false
}

func (t *Table) AddCard(card *Card) {

	for i := range t.Cards {
		if t.Cards[i].UUID == card.UUID {
			return
		}
	}
	t.Cards = append(t.Cards, card)
	t.updateSize()
}

func (t *Table) updateSize() {
	t.Size = len(t.Cards)
}
