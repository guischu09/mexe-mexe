package engine

import "fmt"

type AvailablePlay string

const (
	PLAY_MELD AvailablePlay = "PLAY_MELD"
	DRAW_CARD AvailablePlay = "DRAW_CARD"
	QUIT      AvailablePlay = "QUIT"
)

type Play interface {
	IsValid() bool
	Make() bool
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
		Points: 0,
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

func (p *Player) PlayTurn() bool {

	for {
		play := GetPlay()
		if play.IsValid() {
			if !play.Make() {
				return false
			} else {
				return true
			}
		} else {
			fmt.Println("Invalid play! Please try again.")
		}
	}
}

func GetPlay() Play {
	userInput := GetUserInput()
	return ParseInput(userInput)
}

func GetUserInput() string {
	return fmt.Sprintln("Not implemented")
}

func ParseInput(input string) Play {

	fmt.Println(input)
	fmt.Println("Not implemented")
	return nil

}
