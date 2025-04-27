package engine

import "fmt"

type AvailablePlay string

const (
	DRAW_CARD_AND_PLAY_MELD AvailablePlay = "DRAW_CARD_AND_PLAY_MELD"
	PLAY_MELD               AvailablePlay = "PLAY_MELD"
	DRAW_CARD               AvailablePlay = "DRAW_CARD"
	QUIT                    AvailablePlay = "QUIT"
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

func IsValid(turnState *TurnState, play Play) bool {
	switch play.GetName() {

	case DRAW_CARD_AND_PLAY_MELD:
		fmt.Println("Not implemented")
		return false

	case PLAY_MELD:
		fmt.Println("Not implemented")
		return false

	case DRAW_CARD:
		if turnState.HasPlayedMeld {
			fmt.Println("You can't draw a card after playing a meld")
			return false
		}

		if turnState.HasDrawedCard {
			fmt.Println("You can't draw a card twice in a turn")
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

func Make(play Play, deck *GameDeck, table *Table, player *Player) {

	switch play.GetName() {

	case DRAW_CARD_AND_PLAY_MELD:
		fmt.Println("Not implemented")
		return

	case PLAY_MELD:
		fmt.Println("Not implemented")
		return

	case DRAW_CARD:
		card := deck.DrawCard()
		player.Hand.AddCard(card)
		return

	case QUIT:
		return

	default:
		return
	}

}
