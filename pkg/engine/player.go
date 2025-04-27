package engine

import (
	"fmt"
	"time"
)

const MAX_BUFFER_SIZE = 3
const EXPIRATION_TIME = 30 // seconds

var INPUT_MAPPING = map[string]AvailablePlay{
	"q": QUIT,
	"h": SELECT_HAND,
	"t": SELECT_TABLE,
	"p": PLAY_MELD,
	"d": DRAW_CARD,
	"e": END_TURN,
}

type TurnState struct {
	HasDrawedCard bool
	HasPlayedMeld bool
}

func NewTurnState() TurnState {
	return TurnState{
		HasDrawedCard: false,
		HasPlayedMeld: false,
	}
}

func (t *TurnState) Update(hasDrawedCard bool, hasPlayedMeld bool) {
	t.HasDrawedCard = hasDrawedCard
	t.HasPlayedMeld = hasPlayedMeld
}

func (t *TurnState) Print() {
	fmt.Printf("Has drawed card in this turn: %t\r\n", t.HasDrawedCard)
	fmt.Printf("Has played meld in this turn: %t\r\n", t.HasPlayedMeld)
}

type Player struct {
	Name   string
	Hand   Hand
	Points uint32
}

func NewPlayer(name string, hand Hand, points uint32) Player {
	return Player{
		Name:   name,
		Hand:   hand,
		Points: points,
	}
}

func (p *Player) Print() {
	fmt.Println(p.Name)
}

func (p *Player) PrintHand() {
	p.Hand.Print()
}

func (p *Player) UpdatePoints(points uint32) {
	p.Points = points
}

func (p *Player) PlayTurn(deck *Deck, table *Table, inputProvider InputProvider, outputProvider OutputProvider) AvailablePlay {

	fmt.Printf("%s's turn.\r\n", p.Name)
	turnState := NewTurnState()

	for {
		play := inputProvider.GetPlay()
		if IsValid(&turnState, play, outputProvider) {
			Make(play, deck, table, p, outputProvider)

			if play.GetName() == DRAW_CARD {
				turnState.Update(true, false)
				fmt.Println("You can now play a meld or end your turn.")
				continue
			}
			return play.GetName()

		} else {
			fmt.Println("Invalid play! Please try again.")
		}
	}
}

type ClearTimerCC struct {
	Timeout time.Time
}

func NewClearTimerCC() ClearTimerCC {
	return ClearTimerCC{
		Timeout: time.Now(),
	}
}
func (c *ClearTimerCC) Reset() {
	c.Timeout = time.Now()
}
func (c *ClearTimerCC) IsExpired(expirationTime time.Duration) bool {
	return time.Since(c.Timeout) > expirationTime*time.Millisecond
}

type CmdTimer struct {
	Timeout time.Time
}

func NewCmdTimer() CmdTimer {
	return CmdTimer{
		Timeout: time.Now(),
	}
}
func (c *CmdTimer) Reset() {
	c.Timeout = time.Now()
}
func (c *CmdTimer) IsExpired(expirationTime time.Duration) bool {
	return time.Since(c.Timeout) > expirationTime*time.Second
}
