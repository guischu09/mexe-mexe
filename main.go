package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

const NO_SHUFFLE_SEED uint64 = 0
const UNIQUE_SHUFFLE_SEED uint64 = 1

type CardSuit string

const SPADE CardSuit = "SPADE"
const CLUB CardSuit = "CLUB"
const HEART CardSuit = "HEART"
const DIAMOND CardSuit = "DIAMOND"

type CardColor string

const RED CardColor = "RED"
const BLACK CardColor = "BLACK"

type CardName string

const TWO string = "TWO"
const THREE string = "THREE"
const FOUR string = "FOUR"
const FIVE string = "FIVE"
const SIX string = "SIX"
const SEVEN string = "SEVEN"
const EIGHT string = "EIGHT"
const NINE string = "NINE"
const TEN string = "TEN"
const JACK string = "JACK"
const QUEEN string = "QUEEN"
const KING string = "KING"
const ACE string = "ACE"

type CardValue uint64

const TWO_VALUE CardValue = 2
const THREE_VALUE CardValue = 3
const FOUR_VALUE CardValue = 4
const FIVE_VALUE CardValue = 5
const SIX_VALUE CardValue = 6
const SEVEN_VALUE CardValue = 7
const EIGHT_VALUE CardValue = 8
const NINE_VALUE CardValue = 9
const TEN_VALUE CardValue = 10
const JACK_VALUE CardValue = 11
const QUEEN_VALUE CardValue = 12
const KING_VALUE CardValue = 13
const ACE_VALUE CardValue = 14

type CardSymbol string

const TWO_SPADE_SYMBOL CardSymbol = "🂢"
const THREE_SPADE_SYMBOL CardSymbol = "🂣"
const FOUR_SPADE_SYMBOL CardSymbol = "🂤"
const FIVE_SPADE_SYMBOL CardSymbol = "🂥"
const SIX_SPADE_SYMBOL CardSymbol = "🂦"
const SEVEN_SPADE_SYMBOL CardSymbol = "🂧"
const EIGHT_SPADE_SYMBOL CardSymbol = "🂨"
const NINE_SPADE_SYMBOL CardSymbol = "🂩"
const TEN_SPADE_SYMBOL CardSymbol = "🂪"
const JACK_SPADE_SYMBOL CardSymbol = "🂫"
const QUEEN_SPADE_SYMBOL CardSymbol = "🂭"
const KING_SPADE_SYMBOL CardSymbol = "🂮"
const ACE_SPADE_SYMBOL CardSymbol = "🂡"

const TWO_DIAMOND_SYMBOL CardSymbol = "🃂"
const THREE_DIAMOND_SYMBOL CardSymbol = "🃃"
const FOUR_DIAMOND_SYMBOL CardSymbol = "🃄"
const FIVE_DIAMOND_SYMBOL CardSymbol = "🃅"
const SIX_DIAMOND_SYMBOL CardSymbol = "🃆"
const SEVEN_DIAMOND_SYMBOL CardSymbol = "🃇"
const EIGHT_DIAMOND_SYMBOL CardSymbol = "🃈"
const NINE_DIAMOND_SYMBOL CardSymbol = "🃉"
const TEN_DIAMOND_SYMBOL CardSymbol = "🃊"
const JACK_DIAMOND_SYMBOL CardSymbol = "🃋"
const QUEEN_DIAMOND_SYMBOL CardSymbol = "🃍"
const KING_DIAMOND_SYMBOL CardSymbol = "🃎"
const ACE_DIAMOND_SYMBOL CardSymbol = "🃁"

const TWO_HEART_SYMBOL CardSymbol = "🂲"
const THREE_HEART_SYMBOL CardSymbol = "🂳"
const FOUR_HEART_SYMBOL CardSymbol = "🂴"
const FIVE_HEART_SYMBOL CardSymbol = "🂵"
const SIX_HEART_SYMBOL CardSymbol = "🂶"
const SEVEN_HEART_SYMBOL CardSymbol = "🂷"
const EIGHT_HEART_SYMBOL CardSymbol = "🂸"
const NINE_HEART_SYMBOL CardSymbol = "🂹"
const TEN_HEART_SYMBOL CardSymbol = "🂺"
const JACK_HEART_SYMBOL CardSymbol = "🂻"
const QUEEN_HEART_SYMBOL CardSymbol = "🂽"
const KING_HEART_SYMBOL CardSymbol = "🂾"
const ACE_HEART_SYMBOL CardSymbol = "🂱"

const TWO_CLUB_SYMBOL CardSymbol = "🃒"
const THREE_CLUB_SYMBOL CardSymbol = "🃓"
const FOUR_CLUB_SYMBOL CardSymbol = "🃔"
const FIVE_CLUB_SYMBOL CardSymbol = "🃕"
const SIX_CLUB_SYMBOL CardSymbol = "🃖"
const SEVEN_CLUB_SYMBOL CardSymbol = "🃗"
const EIGHT_CLUB_SYMBOL CardSymbol = "🃘"
const NINE_CLUB_SYMBOL CardSymbol = "🃙"
const TEN_CLUB_SYMBOL CardSymbol = "🃚"
const JACK_CLUB_SYMBOL CardSymbol = "🃛"
const QUEEN_CLUB_SYMBOL CardSymbol = "🃝"
const KING_CLUB_SYMBOL CardSymbol = "🃞"
const ACE_CLUB_SYMBOL CardSymbol = "🃑"

type Card struct {
	Name   string
	Suit   CardSuit
	Value  CardValue
	Symbol CardSymbol
	Color  CardColor
}

func (c *Card) Print() {
	fmt.Println(c.Name + " " + string(c.Symbol))
}

func NewCard(name string, suit CardSuit, value CardValue, symbol CardSymbol, color CardColor) (Card, error) {
	// Fix the condition for hearts and diamonds
	if (suit == HEART || suit == DIAMOND) && color != RED {
		return Card{}, errors.New("Cannot create card with suit " + string(suit) + " and color " + string(color))
	}

	// Fix the condition for clubs and spades
	if (suit == CLUB || suit == SPADE) && color != BLACK {
		return Card{}, errors.New("Cannot create card with suit " + string(suit) + " and color " + string(color))
	}

	fullCardName := strings.ToLower(name + " of " + string(suit) + "s")

	newCard := Card{
		Name:   fullCardName,
		Suit:   suit,
		Value:  value,
		Symbol: symbol,
		Color:  color,
	}

	return newCard, nil
}

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

func (g *GameDeck) Contains(card Card) bool {

	for i := 0; i < len(g.Cards); i++ {
		if g.Cards[i].Name == card.Name {
			return true
		}
	}
	return false
}

func (g *GameDeck) DrawCard() Card {
	card := g.Cards[0]
	g.Cards = g.Cards[1:]
	g.updateSize()
	return card
}

func main() {
	gameDeck := NewGameDeck(UNIQUE_SHUFFLE_SEED)
	gameDeck.PrintSize()
	card := gameDeck.DrawCard()

	gameDeck.PrintSize()
	card.Print()

	// fmt.Println(gameDeck.Size)
	// for i := range gameDeck.Cards {
	// fmt.Println(gameDeck.Cards[i].Name + " " + string(gameDeck.Cards[i].Symbol))
	// }

	// newCard, _ := NewCard(TWO, SPADE, TWO_VALUE, TWO_SPADE_SYMBOL, BLACK)
	// fmt.Println(newCard.Name)
	// fmt.Println(newCard.Suit)
	// fmt.Println(newCard.Value)
	// fmt.Println(newCard.Symbol)
	// fmt.Println(newCard.Color)
}
