package engine

import (
	"errors"
	"fmt"
	"strings"
)

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
	UUID   uint8
}

func PrintCards(cards []Card) {
	printHand := ""
	for i := 0; i < len(cards); i++ {
		printHand += string(cards[i].Symbol) + " "
	}
	fmt.Println(printHand)
}

func (c *Card) Print() {
	fmt.Println(c.Symbol)
}

func NewCard(name string, suit CardSuit, value CardValue, symbol CardSymbol, color CardColor) (*Card, error) {
	// Fix the condition for hearts and diamonds
	if (suit == HEART || suit == DIAMOND) && color != RED {
		return &Card{}, errors.New("Cannot create card with suit " + string(suit) + " and color " + string(color))
	}

	// Fix the condition for clubs and spades
	if (suit == CLUB || suit == SPADE) && color != BLACK {
		return &Card{}, errors.New("Cannot create card with suit " + string(suit) + " and color " + string(color))
	}

	fullCardName := strings.ToLower(name + " of " + string(suit) + "s" + " " + string(symbol))

	newCard := Card{
		Name:   fullCardName,
		Suit:   suit,
		Value:  value,
		Symbol: symbol,
		Color:  color,
	}

	return &newCard, nil
}
