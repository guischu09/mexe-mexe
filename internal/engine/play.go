package engine

import "log"

type AvailablePlay string

const (
	PLAY_MELD    AvailablePlay = "PLAY_MELD"
	DRAW_CARD    AvailablePlay = "DRAW_CARD"
	QUIT         AvailablePlay = "QUIT"
	END_TURN     AvailablePlay = "END_TURN"
	SELECT_HAND  AvailablePlay = "SELECT_HAND"
	SELECT_TABLE AvailablePlay = "SELECT_TABLE"
)

type Play interface {
	GetName() AvailablePlay
	GetCards() []*Card
}

type QuitPlay struct {
	Command string
}

func (q QuitPlay) GetName() AvailablePlay {
	return QUIT
}
func (p QuitPlay) GetCards() []*Card {
	log.Println("No cards used for this play")
	return nil
}

func NewQuitPlay(command string) QuitPlay {
	return QuitPlay{
		Command: command,
	}
}

type EndTurnPlay struct {
	Command string
}

func (p EndTurnPlay) GetName() AvailablePlay {
	return END_TURN
}

func (p EndTurnPlay) GetCards() []*Card {
	log.Println("No cards used for this play")
	return nil
}
func NewEndTurnPlay(command string) EndTurnPlay {
	return EndTurnPlay{
		Command: command,
	}
}

type MeldPlay struct {
	Command string
	Cards   []*Card
}

func (m MeldPlay) GetName() AvailablePlay {
	return PLAY_MELD
}

func (m MeldPlay) GetCards() []*Card {
	return m.Cards
}

func NewMeldPlay(command string, cards []*Card) MeldPlay {
	return MeldPlay{
		Command: command,
		Cards:   cards,
	}
}

type DrawCardPlay struct {
	Command string
}

func (d DrawCardPlay) GetName() AvailablePlay {
	return DRAW_CARD
}
func (p DrawCardPlay) GetCards() []*Card {
	log.Println("No cards used for this play")
	return nil
}

func NewDrawCardPlay(command string) DrawCardPlay {
	return DrawCardPlay{
		Command: command,
	}
}

func IsValid(turnState *TurnState, play Play, outputProvider OutputProvider) bool {
	switch play.GetName() {

	case PLAY_MELD:
		return true

	case END_TURN:

		if turnState.HasPlayedMeld || turnState.HasDrawedCard {
			return true
		} else {
			outputProvider.Write("message", "You must play a meld or draw a card before ending the turn.")
			return false
		}

	case DRAW_CARD:
		if turnState.HasPlayedMeld {
			outputProvider.Write("message", "You can't draw a card after playing a meld.")
			return false
		}

		if turnState.HasDrawedCard {
			outputProvider.Write("message", "You can't draw a card twice in a turn.")
			return false

		} else {
			return true
		}

	case QUIT:
		return true

	default:
		return false
	}
}

func MakePlay(play Play, deck *Deck, table *Table, player *Player, outputProvider OutputProvider) {

	switch play.GetName() {

	case PLAY_MELD:
		for _, card := range play.GetCards() {
			if player.Hand.Contains(card) {
				player.Hand.RemoveCard(card)
				table.AddCard(card)
			}
		}
		return

	case DRAW_CARD:
		card := deck.DrawCard()
		card.Print()
		player.Hand.AddCard(card)
		player.Hand.Print()
		return

	case END_TURN:
		outputProvider.Write("message", "Passing turn")
		return

	case QUIT:
		return

	default:
		return
	}

}
