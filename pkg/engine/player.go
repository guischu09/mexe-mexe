package engine

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const MAX_BUFFER_SIZE = 3

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

func (p *Player) PlayTurn(deck *GameDeck, table *Table) AvailablePlay {

	fmt.Printf("%s's turn.\r\n", p.Name)
	turnState := NewTurnState()

	for {
		play := GetPlayTerminal()
		if IsValid(&turnState, play) {
			Make(play, deck, table, p)
			return play.GetName()

		} else {
			fmt.Println("Invalid play! Please try again.")
		}
	}
}

func GetPlayTerminal() Play {
	for {
		userInput := GetUserInput()
		play := ParseInput(userInput)
		if play == nil {
			continue
		}
		return play
	}
}

func GetUserInput() string {

	// Get keyboard input from stdin
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(0, oldState)

	// var gameCmd []string

	for {
		// Read from stdin into buffer
		buffer := make([]byte, MAX_BUFFER_SIZE)
		_, err = os.Stdin.Read(buffer[:1])
		if err != nil {
			continue
		}

		detectedKey := strings.ToLower(string(buffer[0]))
		fmt.Printf("Key pressed: '%s'\r\n", detectedKey)

		if detectedKey == "q" || detectedKey == "d" {
			return detectedKey
		}

	}
}

func ParseInput(input string) Play {

	switch input {
	case "q":
		return NewQuitPlay(input)
	case "d":
		return NewDrawCardPlay(input)
	// case "m":
	// return NewMeldPlay(input)
	default:
		return nil
	}

}
