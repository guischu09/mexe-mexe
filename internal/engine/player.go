package engine

import (
	"fmt"
	"log"
)

const MAX_BUFFER_SIZE = 3
const EXPIRATION_TIME = 30 // seconds

var INPUT_MAPPING = map[string]AvailablePlay{
	"q": QUIT,
	"m": PLAY_MELD,
	"d": DRAW_CARD,
	"e": END_TURN,
}

type TurnState struct {
	HasDrawedCard bool
	HasPlayedMeld bool
	PlayerUUID    string
}

func NewTurnState(playerUUID string) *TurnState {
	return &TurnState{
		HasDrawedCard: false,
		HasPlayedMeld: false,
		PlayerUUID:    playerUUID,
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
	Hand   *Hand
	Points uint32
	UUID   string
}

func NewPlayer(name string, hand *Hand, uuid string, points uint32) Player {
	return Player{
		Name:   name,
		Hand:   hand,
		Points: points,
		UUID:   uuid,
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

	turnState := NewTurnState(p.UUID)

	for {
		log.Print("player :: !> DEBUG: Turn state: turnState.HasDrawedCard: ", turnState.HasDrawedCard)
		log.Print("player :: !> DEBUG: Turn state: turnState.HasPlayedMeld: ", turnState.HasPlayedMeld)
		log.Print("player :: !> DEBUG: Turn state: turnState.PlayerUUID: ", turnState.PlayerUUID)

		outputProvider.SendState(table, *p.Hand, *turnState)
		play := inputProvider.GetPlay(table, *p.Hand, p.Name, *turnState)
		log.Print("player :: !> Got Play: ", play.GetName())
		if IsValid(turnState, play, outputProvider) {
			log.Print("player :: !> Play is valid")

			MakePlay(play, deck, table, p, outputProvider)

			if play.GetName() == DRAW_CARD {
				turnState.Update(true, false)
				continue
			}

			if play.GetName() == PLAY_MELD {
				turnState.Update(false, true)
				continue
			}

			return play.GetName()

		} else {
			log.Print("player :: !> Play is valid")
			fmt.Println("Invalid play! Please try again.")
		}
	}
}
