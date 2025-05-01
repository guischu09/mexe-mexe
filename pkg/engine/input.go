package engine

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type InputProvider interface {
	GetPlay(*Table) Play
}

type TerminalInputProvider struct{}

func (t *TerminalInputProvider) GetPlay(table *Table) Play {
	for {
		userInput := GetUserInput(table)
		play := ParseInput(userInput)
		if play == nil {
			continue
		}
		return play
	}
}

func GetUserInput(table *Table) string {

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
				fmt.Printf("Continue playing...\r\n")
				continue
			}
		}

		if detectedKey == "m" {
			fmt.Printf("Select your cards and play a meld.\r\n")
			MeldDisplayInput(table)
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

func MeldDisplayInput(table *Table) {

	width, height, err := term.GetSize(0)
	if err != nil {
		return
	}
	fmt.Printf("width:%d height:%d\r\n", width, height)

	table.Print(width, height)

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

type NetworkInputProvider struct {
	// conn websocket.Conn
}

func (n *NetworkInputProvider) GetPlay() Play {
	return nil
}
