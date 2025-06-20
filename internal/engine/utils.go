package engine

import "sort"

func SortCardsByValue(cards []Card) {
	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			if cards[i].Value > cards[j].Value {
				cards[i], cards[j] = cards[j], cards[i]
			}
		}
	}
}

func SortCardsBySuit(cards []Card) {
	suitOrder := map[CardSuit]int{
		SPADE:   0,
		CLUB:    1,
		HEART:   2,
		DIAMOND: 3,
	}

	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			if suitOrder[cards[i].Suit] > suitOrder[cards[j].Suit] {
				cards[i], cards[j] = cards[j], cards[i]
			}
		}
	}
}

func SortCardsBySuitAndValue(cards []*Card) {
	suitOrder := map[CardSuit]int{
		SPADE: 0, CLUB: 1, HEART: 2, DIAMOND: 3,
	}

	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Suit != cards[j].Suit {
			return suitOrder[cards[i].Suit] < suitOrder[cards[j].Suit]
		}
		return cards[i].Value < cards[j].Value
	})
}

func SortHandBySuitAndValue(hand *Hand) {
	suitOrder := map[CardSuit]int{
		SPADE: 0, CLUB: 1, HEART: 2, DIAMOND: 3,
	}

	sort.Slice(hand.Cards, func(i, j int) bool {
		if hand.Cards[i].Suit != hand.Cards[j].Suit {
			return suitOrder[hand.Cards[i].Suit] < suitOrder[hand.Cards[j].Suit]
		}
		return hand.Cards[i].Value < hand.Cards[j].Value
	})
}
