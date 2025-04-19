package engine

import "fmt"

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
