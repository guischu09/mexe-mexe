package engine

func SortCardsByValue(cards []*Card) {
	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			if cards[i].Value > cards[j].Value {
				cards[i], cards[j] = cards[j], cards[i]
			}
		}
	}
}

func SortCardsBySuit(cards []*Card) {
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
		SPADE:   0,
		CLUB:    1,
		HEART:   2,
		DIAMOND: 3,
	}

	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			if suitOrder[cards[i].Suit] > suitOrder[cards[j].Suit] {
				cards[i], cards[j] = cards[j], cards[i]
			} else if suitOrder[cards[i].Suit] == suitOrder[cards[j].Suit] {
				if cards[i].Value > cards[j].Value {
					cards[i], cards[j] = cards[j], cards[i]
				}
			}
		}
	}
}
