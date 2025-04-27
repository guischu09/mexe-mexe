package engine

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
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

func (p *Player) PlayTurn(deck *GameDeck, table *Table) AvailablePlay {

	fmt.Printf("%s's turn.\r\n", p.Name)
	turnState := NewTurnState()

	for {
		play := GetPlayTerminal()
		if IsValid(&turnState, play) {
			Make(play, deck, table, p)

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

func GetUserInput() string {

	// Get keyboard input from stdin
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(0, oldState)

	var gameCmd []string

	cmdTimer := NewCmdTimer()
	ccTimer := NewClearTimerCC()

	for {
		// Read from stdin into buffer
		buffer := make([]byte, MAX_BUFFER_SIZE)
		_, err = os.Stdin.Read(buffer[:1])
		if err != nil {
			continue
		}

		detectedKey := strings.ToLower(string(buffer[0]))
		fmt.Printf("Key pressed: '%s'\r\n", detectedKey)

		// Handle quit
		if detectedKey == "q" {
			fmt.Printf("Are you sure you want to quit? (y/n)\r\n")
			_, err = os.Stdin.Read(buffer[:1])
			if err != nil {
				continue
			}
			if strings.ToLower(string(buffer[0])) == "y" || strings.ToLower(string(buffer[0])) == "q" {
				return "q"
			} else {
				continue
			}
		}

		if detectedKey == "d" {
			return detectedKey
		}

		if detectedKey == "e" {
			return detectedKey
		}

		// clear input
		if detectedKey == "c" {
			if len(gameCmd) > 0 {
				gameCmd = gameCmd[:len(gameCmd)-1]
			}
			fmt.Printf("Game command: %s\r\n", gameCmd)
			ccTimer.Reset()
			cmdTimer.Reset()
			continue
		}

		// append input to gameCmd
		gameCmd = append(gameCmd, detectedKey)
		fmt.Printf("Game command: %s\r\n", gameCmd)

		// Clear gameCmd when timeout
		if cmdTimer.IsExpired(EXPIRATION_TIME) {
			gameCmd = []string{}
			cmdTimer.Reset()
			fmt.Printf("Command timeout expired!\r\n")
			fmt.Printf("Game command: %s\r\n", gameCmd)
		}
	}
}

func ParseInput(input string) Play {

	switch input {
	case "q":
		return NewQuitPlay(input)
	case "d":
		return NewDrawCardPlay(input)
	case "e":
		return NewPassPlay(input)
	// case "m":
	// return NewMeldPlay(input)
	default:
		return nil
	}

}
