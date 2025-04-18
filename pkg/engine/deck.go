package engine

import (
	"fmt"
	"math/rand/v2"
	"time"
)

const NO_SHUFFLE_SEED uint64 = 0
const UNIQUE_SHUFFLE_SEED uint64 = 1

type GameDeck struct {
	Cards []Card
	Size  int
	Seed  uint64
}

func NewGameDeck(seed uint64) GameDeck {
	var cards []Card

	twoOfSpades, _ := NewCard(TWO, SPADE, TWO_VALUE, TWO_SPADE_SYMBOL, BLACK)
	threeOfSpades, _ := NewCard(THREE, SPADE, THREE_VALUE, THREE_SPADE_SYMBOL, BLACK)
	fourOfSpades, _ := NewCard(FOUR, SPADE, FOUR_VALUE, FOUR_SPADE_SYMBOL, BLACK)
	fiveOfSpades, _ := NewCard(FIVE, SPADE, FIVE_VALUE, FIVE_SPADE_SYMBOL, BLACK)
	sixOfSpades, _ := NewCard(SIX, SPADE, SIX_VALUE, SIX_SPADE_SYMBOL, BLACK)
	sevenOfSpades, _ := NewCard(SEVEN, SPADE, SEVEN_VALUE, SEVEN_SPADE_SYMBOL, BLACK)
	eightOfSpades, _ := NewCard(EIGHT, SPADE, EIGHT_VALUE, EIGHT_SPADE_SYMBOL, BLACK)
	nineOfSpades, _ := NewCard(NINE, SPADE, NINE_VALUE, NINE_SPADE_SYMBOL, BLACK)
	tenOfSpades, _ := NewCard(TEN, SPADE, TEN_VALUE, TEN_SPADE_SYMBOL, BLACK)
	jackOfSpades, _ := NewCard(JACK, SPADE, JACK_VALUE, JACK_SPADE_SYMBOL, BLACK)
	queenOfSpades, _ := NewCard(QUEEN, SPADE, QUEEN_VALUE, QUEEN_SPADE_SYMBOL, BLACK)
	kingOfSpades, _ := NewCard(KING, SPADE, KING_VALUE, KING_SPADE_SYMBOL, BLACK)
	aceOfSpades, _ := NewCard(ACE, SPADE, ACE_VALUE, ACE_SPADE_SYMBOL, BLACK)

	cards = append(cards, twoOfSpades, threeOfSpades, fourOfSpades, fiveOfSpades, sixOfSpades, sevenOfSpades, eightOfSpades, nineOfSpades, tenOfSpades, jackOfSpades, queenOfSpades, kingOfSpades, aceOfSpades)
	cards = append(cards, twoOfSpades, threeOfSpades, fourOfSpades, fiveOfSpades, sixOfSpades, sevenOfSpades, eightOfSpades, nineOfSpades, tenOfSpades, jackOfSpades, queenOfSpades, kingOfSpades, aceOfSpades)

	twoOfDiamonds, _ := NewCard(TWO, DIAMOND, TWO_VALUE, TWO_DIAMOND_SYMBOL, RED)
	threeOfDiamonds, _ := NewCard(THREE, DIAMOND, THREE_VALUE, THREE_DIAMOND_SYMBOL, RED)
	fourOfDiamonds, _ := NewCard(FOUR, DIAMOND, FOUR_VALUE, FOUR_DIAMOND_SYMBOL, RED)
	fiveOfDiamonds, _ := NewCard(FIVE, DIAMOND, FIVE_VALUE, FIVE_DIAMOND_SYMBOL, RED)
	sixOfDiamonds, _ := NewCard(SIX, DIAMOND, SIX_VALUE, SIX_DIAMOND_SYMBOL, RED)
	sevenOfDiamonds, _ := NewCard(SEVEN, DIAMOND, SEVEN_VALUE, SEVEN_DIAMOND_SYMBOL, RED)
	eightOfDiamonds, _ := NewCard(EIGHT, DIAMOND, EIGHT_VALUE, EIGHT_DIAMOND_SYMBOL, RED)
	nineOfDiamonds, _ := NewCard(NINE, DIAMOND, NINE_VALUE, NINE_DIAMOND_SYMBOL, RED)
	tenOfDiamonds, _ := NewCard(TEN, DIAMOND, TEN_VALUE, TEN_DIAMOND_SYMBOL, RED)
	jackOfDiamonds, _ := NewCard(JACK, DIAMOND, JACK_VALUE, JACK_DIAMOND_SYMBOL, RED)
	queenOfDiamonds, _ := NewCard(QUEEN, DIAMOND, QUEEN_VALUE, QUEEN_DIAMOND_SYMBOL, RED)
	kingOfDiamonds, _ := NewCard(KING, DIAMOND, KING_VALUE, KING_DIAMOND_SYMBOL, RED)
	aceOfDiamonds, _ := NewCard(ACE, DIAMOND, ACE_VALUE, ACE_DIAMOND_SYMBOL, RED)

	cards = append(cards, twoOfDiamonds, threeOfDiamonds, fourOfDiamonds, fiveOfDiamonds, sixOfDiamonds, sevenOfDiamonds, eightOfDiamonds, nineOfDiamonds, tenOfDiamonds, jackOfDiamonds, queenOfDiamonds, kingOfDiamonds, aceOfDiamonds)
	cards = append(cards, twoOfDiamonds, threeOfDiamonds, fourOfDiamonds, fiveOfDiamonds, sixOfDiamonds, sevenOfDiamonds, eightOfDiamonds, nineOfDiamonds, tenOfDiamonds, jackOfDiamonds, queenOfDiamonds, kingOfDiamonds, aceOfDiamonds)

	twoOfHearts, _ := NewCard(TWO, HEART, TWO_VALUE, TWO_HEART_SYMBOL, RED)
	threeOfHearts, _ := NewCard(THREE, HEART, THREE_VALUE, THREE_HEART_SYMBOL, RED)
	fourOfHearts, _ := NewCard(FOUR, HEART, FOUR_VALUE, FOUR_HEART_SYMBOL, RED)
	fiveOfHearts, _ := NewCard(FIVE, HEART, FIVE_VALUE, FIVE_HEART_SYMBOL, RED)
	sixOfHearts, _ := NewCard(SIX, HEART, SIX_VALUE, SIX_HEART_SYMBOL, RED)
	sevenOfHearts, _ := NewCard(SEVEN, HEART, SEVEN_VALUE, SEVEN_HEART_SYMBOL, RED)
	eightOfHearts, _ := NewCard(EIGHT, HEART, EIGHT_VALUE, EIGHT_HEART_SYMBOL, RED)
	nineOfHearts, _ := NewCard(NINE, HEART, NINE_VALUE, NINE_HEART_SYMBOL, RED)
	tenOfHearts, _ := NewCard(TEN, HEART, TEN_VALUE, TEN_HEART_SYMBOL, RED)
	jackOfHearts, _ := NewCard(JACK, HEART, JACK_VALUE, JACK_HEART_SYMBOL, RED)
	queenOfHearts, _ := NewCard(QUEEN, HEART, QUEEN_VALUE, QUEEN_HEART_SYMBOL, RED)
	kingOfHearts, _ := NewCard(KING, HEART, KING_VALUE, KING_HEART_SYMBOL, RED)
	aceOfHearts, _ := NewCard(ACE, HEART, ACE_VALUE, ACE_HEART_SYMBOL, RED)

	cards = append(cards, twoOfHearts, threeOfHearts, fourOfHearts, fiveOfHearts, sixOfHearts, sevenOfHearts, eightOfHearts, nineOfHearts, tenOfHearts, jackOfHearts, queenOfHearts, kingOfHearts, aceOfHearts)
	cards = append(cards, twoOfHearts, threeOfHearts, fourOfHearts, fiveOfHearts, sixOfHearts, sevenOfHearts, eightOfHearts, nineOfHearts, tenOfHearts, jackOfHearts, queenOfHearts, kingOfHearts, aceOfHearts)

	twoOfClubs, _ := NewCard(TWO, CLUB, TWO_VALUE, TWO_CLUB_SYMBOL, BLACK)
	threeOfClubs, _ := NewCard(THREE, CLUB, THREE_VALUE, THREE_CLUB_SYMBOL, BLACK)
	fourOfClubs, _ := NewCard(FOUR, CLUB, FOUR_VALUE, FOUR_CLUB_SYMBOL, BLACK)
	fiveOfClubs, _ := NewCard(FIVE, CLUB, FIVE_VALUE, FIVE_CLUB_SYMBOL, BLACK)
	sixOfClubs, _ := NewCard(SIX, CLUB, SIX_VALUE, SIX_CLUB_SYMBOL, BLACK)
	sevenOfClubs, _ := NewCard(SEVEN, CLUB, SEVEN_VALUE, SEVEN_CLUB_SYMBOL, BLACK)
	eightOfClubs, _ := NewCard(EIGHT, CLUB, EIGHT_VALUE, EIGHT_CLUB_SYMBOL, BLACK)
	nineOfClubs, _ := NewCard(NINE, CLUB, NINE_VALUE, NINE_CLUB_SYMBOL, BLACK)
	tenOfClubs, _ := NewCard(TEN, CLUB, TEN_VALUE, TEN_CLUB_SYMBOL, BLACK)
	jackOfClubs, _ := NewCard(JACK, CLUB, JACK_VALUE, JACK_CLUB_SYMBOL, BLACK)
	queenOfClubs, _ := NewCard(QUEEN, CLUB, QUEEN_VALUE, QUEEN_CLUB_SYMBOL, BLACK)
	kingOfClubs, _ := NewCard(KING, CLUB, KING_VALUE, KING_CLUB_SYMBOL, BLACK)
	aceOfClubs, _ := NewCard(ACE, CLUB, ACE_VALUE, ACE_CLUB_SYMBOL, BLACK)

	cards = append(cards, twoOfClubs, threeOfClubs, fourOfClubs, fiveOfClubs, sixOfClubs, sevenOfClubs, eightOfClubs, nineOfClubs, tenOfClubs, jackOfClubs, queenOfClubs, kingOfClubs, aceOfClubs)
	cards = append(cards, twoOfClubs, threeOfClubs, fourOfClubs, fiveOfClubs, sixOfClubs, sevenOfClubs, eightOfClubs, nineOfClubs, tenOfClubs, jackOfClubs, queenOfClubs, kingOfClubs, aceOfClubs)

	if seed == NO_SHUFFLE_SEED {
		gameDeck := GameDeck{
			Cards: cards,
			Size:  len(cards),
			Seed:  seed,
		}
		return gameDeck
	}

	if seed == UNIQUE_SHUFFLE_SEED {
		seedTime := uint64(time.Now().UnixNano())

		// Use seedTime for PCG initialization, not seed
		pcgSource := rand.NewPCG(seedTime, seedTime)
		rng := rand.New(pcgSource)

		rng.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
		gameDeck := GameDeck{
			Cards: cards,
			Size:  len(cards),
			Seed:  seedTime,
		}
		return gameDeck
	}

	pcgSource := rand.NewPCG(seed, seed)
	rng := rand.New(pcgSource)
	rng.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	gameDeck := GameDeck{
		Cards: cards,
		Size:  len(cards),
		Seed:  seed,
	}

	return gameDeck
}

func (g *GameDeck) Print() {
	for i := 0; i < len(g.Cards); i++ {
		fmt.Println(g.Cards[i].Name + " " + string(g.Cards[i].Symbol))
	}
	fmt.Println("Size: ", g.Size)
}

func (g *GameDeck) PrintSize() {
	fmt.Println("Size: ", g.Size)
}

func (g *GameDeck) updateSize() {
	g.Size = len(g.Cards)
}

// func binary_search

func (g *GameDeck) Contains(card Card, option string) bool {
	switch option {
	case "binary_search":
		fmt.Println("binary_search not implemented")
		return false
	default:
		for i := 0; i < len(g.Cards); i++ {
			if g.Cards[i].Name == card.Name {
				return true
			}
		}
		return false
	}
}

func (g *GameDeck) DrawCard() Card {
	card := g.Cards[0]
	g.Cards = g.Cards[1:]
	g.updateSize()
	return card
}

func (g *GameDeck) RemoveCard(card Card, option string) bool {
	switch option {
	case "binary_search":
		panic("binary_search not implemented")
	default:
		for i := 0; i < len(g.Cards); i++ {
			if g.Cards[i].Name == card.Name {
				g.Cards = append(g.Cards[:i], g.Cards[i+1:]...)
				g.updateSize()
				return true
			}
		}
		fmt.Println("Card not found in the deck")
		return false
	}
}
