package engine

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

const TOTAL_DECK_SIZE uint64 = 104
const NO_SHUFFLE_SEED uint64 = 0
const UNIQUE_SHUFFLE_SEED uint64 = 1

type Deck struct {
	Cards       []*Card
	Size        int
	Seed        uint64
	uuidCounter uint8
	mu          sync.Mutex // for thread safety if needed
}

// Card, CardSuit, CardColor, CardValue, CardSymbol, and all constants should be defined as in your original code.

func (d *Deck) getNextCardUUID() uint8 {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.uuidCounter++
	return d.uuidCounter
}

func (d *Deck) newCard(name string, suit CardSuit, value CardValue, symbol CardSymbol, color CardColor) *Card {
	uuid := d.getNextCardUUID()
	fullCardName := name + " of " + string(suit) + "s " + string(symbol)
	return &Card{
		Name:   fullCardName,
		Suit:   suit,
		Value:  value,
		Symbol: symbol,
		Color:  color,
		UUID:   uuid,
	}
}

func NewDeck(seed uint64) *Deck {
	deck := Deck{
		Cards:       []*Card{},
		Size:        0,
		Seed:        seed,
		uuidCounter: 0,
	}

	// Helper to add two of each card
	add := func(name string, suit CardSuit, value CardValue, symbol CardSymbol, color CardColor) {
		deck.Cards = append(deck.Cards, deck.newCard(name, suit, value, symbol, color))
		deck.Cards = append(deck.Cards, deck.newCard(name, suit, value, symbol, color))
	}

	// Spades
	add(TWO, SPADE, TWO_VALUE, TWO_SPADE_SYMBOL, BLACK)
	add(THREE, SPADE, THREE_VALUE, THREE_SPADE_SYMBOL, BLACK)
	add(FOUR, SPADE, FOUR_VALUE, FOUR_SPADE_SYMBOL, BLACK)
	add(FIVE, SPADE, FIVE_VALUE, FIVE_SPADE_SYMBOL, BLACK)
	add(SIX, SPADE, SIX_VALUE, SIX_SPADE_SYMBOL, BLACK)
	add(SEVEN, SPADE, SEVEN_VALUE, SEVEN_SPADE_SYMBOL, BLACK)
	add(EIGHT, SPADE, EIGHT_VALUE, EIGHT_SPADE_SYMBOL, BLACK)
	add(NINE, SPADE, NINE_VALUE, NINE_SPADE_SYMBOL, BLACK)
	add(TEN, SPADE, TEN_VALUE, TEN_SPADE_SYMBOL, BLACK)
	add(JACK, SPADE, JACK_VALUE, JACK_SPADE_SYMBOL, BLACK)
	add(QUEEN, SPADE, QUEEN_VALUE, QUEEN_SPADE_SYMBOL, BLACK)
	add(KING, SPADE, KING_VALUE, KING_SPADE_SYMBOL, BLACK)
	add(ACE, SPADE, ACE_VALUE, ACE_SPADE_SYMBOL, BLACK)

	// Diamonds
	add(TWO, DIAMOND, TWO_VALUE, TWO_DIAMOND_SYMBOL, RED)
	add(THREE, DIAMOND, THREE_VALUE, THREE_DIAMOND_SYMBOL, RED)
	add(FOUR, DIAMOND, FOUR_VALUE, FOUR_DIAMOND_SYMBOL, RED)
	add(FIVE, DIAMOND, FIVE_VALUE, FIVE_DIAMOND_SYMBOL, RED)
	add(SIX, DIAMOND, SIX_VALUE, SIX_DIAMOND_SYMBOL, RED)
	add(SEVEN, DIAMOND, SEVEN_VALUE, SEVEN_DIAMOND_SYMBOL, RED)
	add(EIGHT, DIAMOND, EIGHT_VALUE, EIGHT_DIAMOND_SYMBOL, RED)
	add(NINE, DIAMOND, NINE_VALUE, NINE_DIAMOND_SYMBOL, RED)
	add(TEN, DIAMOND, TEN_VALUE, TEN_DIAMOND_SYMBOL, RED)
	add(JACK, DIAMOND, JACK_VALUE, JACK_DIAMOND_SYMBOL, RED)
	add(QUEEN, DIAMOND, QUEEN_VALUE, QUEEN_DIAMOND_SYMBOL, RED)
	add(KING, DIAMOND, KING_VALUE, KING_DIAMOND_SYMBOL, RED)
	add(ACE, DIAMOND, ACE_VALUE, ACE_DIAMOND_SYMBOL, RED)

	// Hearts
	add(TWO, HEART, TWO_VALUE, TWO_HEART_SYMBOL, RED)
	add(THREE, HEART, THREE_VALUE, THREE_HEART_SYMBOL, RED)
	add(FOUR, HEART, FOUR_VALUE, FOUR_HEART_SYMBOL, RED)
	add(FIVE, HEART, FIVE_VALUE, FIVE_HEART_SYMBOL, RED)
	add(SIX, HEART, SIX_VALUE, SIX_HEART_SYMBOL, RED)
	add(SEVEN, HEART, SEVEN_VALUE, SEVEN_HEART_SYMBOL, RED)
	add(EIGHT, HEART, EIGHT_VALUE, EIGHT_HEART_SYMBOL, RED)
	add(NINE, HEART, NINE_VALUE, NINE_HEART_SYMBOL, RED)
	add(TEN, HEART, TEN_VALUE, TEN_HEART_SYMBOL, RED)
	add(JACK, HEART, JACK_VALUE, JACK_HEART_SYMBOL, RED)
	add(QUEEN, HEART, QUEEN_VALUE, QUEEN_HEART_SYMBOL, RED)
	add(KING, HEART, KING_VALUE, KING_HEART_SYMBOL, RED)
	add(ACE, HEART, ACE_VALUE, ACE_HEART_SYMBOL, RED)

	// Clubs
	add(TWO, CLUB, TWO_VALUE, TWO_CLUB_SYMBOL, BLACK)
	add(THREE, CLUB, THREE_VALUE, THREE_CLUB_SYMBOL, BLACK)
	add(FOUR, CLUB, FOUR_VALUE, FOUR_CLUB_SYMBOL, BLACK)
	add(FIVE, CLUB, FIVE_VALUE, FIVE_CLUB_SYMBOL, BLACK)
	add(SIX, CLUB, SIX_VALUE, SIX_CLUB_SYMBOL, BLACK)
	add(SEVEN, CLUB, SEVEN_VALUE, SEVEN_CLUB_SYMBOL, BLACK)
	add(EIGHT, CLUB, EIGHT_VALUE, EIGHT_CLUB_SYMBOL, BLACK)
	add(NINE, CLUB, NINE_VALUE, NINE_CLUB_SYMBOL, BLACK)
	add(TEN, CLUB, TEN_VALUE, TEN_CLUB_SYMBOL, BLACK)
	add(JACK, CLUB, JACK_VALUE, JACK_CLUB_SYMBOL, BLACK)
	add(QUEEN, CLUB, QUEEN_VALUE, QUEEN_CLUB_SYMBOL, BLACK)
	add(KING, CLUB, KING_VALUE, KING_CLUB_SYMBOL, BLACK)
	add(ACE, CLUB, ACE_VALUE, ACE_CLUB_SYMBOL, BLACK)

	// Verify we have the correct number of cards
	if len(deck.Cards) != int(TOTAL_DECK_SIZE) {
		panic(fmt.Sprintf("Expected %d cards in deck, got %d", TOTAL_DECK_SIZE, len(deck.Cards)))
	}

	// Shuffle logic
	if seed == NO_SHUFFLE_SEED {
		deck.Size = len(deck.Cards)
		return &deck
	}

	if seed == UNIQUE_SHUFFLE_SEED {
		seedTime := uint64(time.Now().UnixNano())
		pcgSource := rand.NewPCG(seedTime, seedTime)
		rng := rand.New(pcgSource)
		rng.Shuffle(len(deck.Cards), func(i, j int) {
			deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
		})
		deck.Size = len(deck.Cards)
		deck.Seed = seedTime
		return &deck
	}

	pcgSource := rand.NewPCG(seed, seed)
	rng := rand.New(pcgSource)
	rng.Shuffle(len(deck.Cards), func(i, j int) {
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	})
	deck.Size = len(deck.Cards)
	return &deck
}

func (g *Deck) Print() {
	for i := 0; i < len(g.Cards); i++ {
		fmt.Println(g.Cards[i].Name + " " + string(g.Cards[i].Symbol))
	}
	fmt.Println("Size: ", g.Size)
}

func (g *Deck) PrintSize() {
	fmt.Println("Size: ", g.Size)
}

func (g *Deck) updateSize() {
	g.Size = len(g.Cards)
}

func (g *Deck) Contains(card *Card) bool {
	for i := 0; i < len(g.Cards); i++ {
		if g.Cards[i].UUID == card.UUID {
			return true
		}
	}
	return false
}

func (g *Deck) DrawCard() *Card {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.Cards) == 0 {
		return nil
	}
	card := g.Cards[0]
	g.Cards = g.Cards[1:]
	g.updateSize()
	return card
}

func (g *Deck) RemoveCard(card *Card) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i := 0; i < len(g.Cards); i++ {
		if g.Cards[i].UUID == card.UUID {
			g.Cards = append(g.Cards[:i], g.Cards[i+1:]...)
			g.updateSize()
			return true
		}
	}
	fmt.Println("Card not found in the deck")
	return false
}
