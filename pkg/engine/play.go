package engine

import "fmt"

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
}

type QuitPlay struct {
	Command string
}

func (q QuitPlay) GetName() AvailablePlay {
	return QUIT
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

func NewPassPlay(command string) EndTurnPlay {
	return EndTurnPlay{
		Command: command,
	}
}

type MeldPlay struct {
	Command string
}

func NewMeldPlay(command string) MeldPlay {
	return MeldPlay{
		Command: command,
	}
}

type DrawCardPlay struct {
	Command string
}

func (d DrawCardPlay) GetName() AvailablePlay {
	return DRAW_CARD
}

func NewDrawCardPlay(command string) DrawCardPlay {
	return DrawCardPlay{
		Command: command,
	}
}

func IsValid(turnState *TurnState, play Play, outputProvider OutputProvider) bool {
	switch play.GetName() {

	case PLAY_MELD:
		outputProvider.Write(string(PLAY_MELD), "Not implemented")
		return false

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

func Make(play Play, deck *Deck, table *Table, player *Player, outputProvider OutputProvider) {

	switch play.GetName() {

	case PLAY_MELD:
		fmt.Println("Not implemented")
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
