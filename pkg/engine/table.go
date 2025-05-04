package engine

import (
	"fmt"
	"strings"
)

type Table struct {
	PlayedCards []*Card
	Size        int
}

func (t *Table) Print(terminalWidth int, terminalHeight int) {
	// for i := 0; i < t.Size; i++ {
	// fmt.Printf("%d ", t.PlayedCards[i].Value)
	// }

	for i := 0; i < terminalHeight*2; i++ {
		for j := 0; j < terminalWidth; j++ {
			if i == j {
				line := getLine()
				fmt.Print(line)
			}
			// fmt.Println("#")
		}
	}
}

// for i := 0; i < terminalHeight; i++ {
// 	for j := 0; j < terminalWidth; j++ {
// 		if i < t.Size && j < len(t.PlayedCards[i].Symbol) {
// 			fmt.Printf("%c", t.PlayedCards[i].Symbol[j])
// 		} else {
// 			fmt.Printf(" ")
// 		}
// 	}
// 	fmt.Println()
// }

// for i := 0; i < terminalHeight; i++ {

func getLine() string {
	line := strings.Repeat("#", 2)
	return line
}

func getColumn() string {
	return ""
}

// func getColumn(terminalWidth int, terminalHeight int) string {

// }
