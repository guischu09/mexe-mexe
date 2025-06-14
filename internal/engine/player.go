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
	GameEnded     bool
}

func NewTurnState(playerUUID string) *TurnState {
	return &TurnState{
		HasDrawedCard: false,
		HasPlayedMeld: false,
		PlayerUUID:    playerUUID,
		GameEnded:     false,
	}
}

func (t *TurnState) UpdateDrawedCard(hasDrawedCard bool) {
	t.HasDrawedCard = hasDrawedCard
}

func (t *TurnState) UpdatePlayedMeld(hasPlayedMeld bool) {
	t.HasPlayedMeld = hasPlayedMeld
}

func (t *TurnState) UpdateGameEnded(hasGameEnded bool) {
	t.GameEnded = hasGameEnded
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

func GetOutputProviderFromUUID(uuid string, outputProviders []OutputProvider) OutputProvider {

	for _, outputProvider := range outputProviders {
		if outputProvider.GetUUID() == uuid {
			return outputProvider
		}
	}
	log.Panicf("ERROR: Output provider not found for UUID: %s", uuid)
	return EMPTY_WS_OUTPUT_PROVIDER
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

func (p *Player) PlayTurn(deck *Deck, table *Table, inputProvider InputProvider, outputProviders []OutputProvider) AvailablePlay {

	turnState := NewTurnState(p.UUID)
	// When the game starts, the server sends the initial state for each player, so that the client
	// can know what to display. We only send another state to the client when the state changes.
	// This is to prevent the client from displaying the same state multiple times.
	// outputProvider.SendState(*table, *p.Hand, *turnState)

	thisPlayerOutputProvider := GetOutputProviderFromUUID(p.UUID, outputProviders)

	for {
		log.Print("player :: !> DEBUG: Turn state: turnState.HasDrawedCard: ", turnState.HasDrawedCard)
		log.Print("player :: !> DEBUG: Turn state: turnState.HasPlayedMeld: ", turnState.HasPlayedMeld)
		log.Print("player :: !> DEBUG: Turn state: turnState.PlayerUUID: ", turnState.PlayerUUID)

		play := inputProvider.GetPlay(*turnState)

		log.Print("player :: !> Got Play: ", play.GetName())
		if IsValid(turnState, play, thisPlayerOutputProvider) {
			log.Print("player :: !> Play is valid")

			MakePlay(play, deck, table, p)

			if play.GetName() == DRAW_CARD {
				turnState.UpdateDrawedCard(true)
				SendStateToPlayers(outputProviders, *table, *p.Hand, *turnState)
				continue
			}

			if play.GetName() == PLAY_MELD {
				turnState.UpdatePlayedMeld(true)
				SendStateToPlayers(outputProviders, *table, *p.Hand, *turnState)
				continue
			}
			if play.GetName() == QUIT {
				turnState.UpdateGameEnded(true)
				SendStateToPlayers(outputProviders, *table, *p.Hand, *turnState)
				return play.GetName()
			}
			log.Panic("Unreachable state reached!")

		} else {
			fmt.Println("Invalid play! Please try again.")
		}
	}
}
