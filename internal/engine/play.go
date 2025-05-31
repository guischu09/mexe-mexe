package engine

import "log"

type AvailablePlay string

const (
	PLAY_MELD AvailablePlay = "PLAY_MELD"
	DRAW_CARD AvailablePlay = "DRAW_CARD"
	QUIT      AvailablePlay = "QUIT"
	END_TURN  AvailablePlay = "END_TURN"
	// SELECT_HAND  AvailablePlay = "SELECT_HAND" deprecated
	// SELECT_TABLE AvailablePlay = "SELECT_TABLE" deprecated
)

type Play interface {
	GetName() AvailablePlay
	GetCards() []Card
}

type QuitPlay struct {
	Type string `json:"type"`
}

func NewQuitPlay() QuitPlay {
	return QuitPlay{Type: "QUIT"}
}

func (q QuitPlay) GetName() AvailablePlay {
	return QUIT
}
func (p QuitPlay) GetCards() []Card {
	log.Println("No cards used for this play")
	return nil
}

type EndTurnPlay struct {
	Type string `json:"type"`
}

func NewEndTurnPlay() EndTurnPlay {
	return EndTurnPlay{Type: "END_TURN"}
}

func (p EndTurnPlay) GetName() AvailablePlay {
	return END_TURN
}

func (p EndTurnPlay) GetCards() []Card {
	log.Println("No cards used for this play")
	return nil
}

type MeldPlay struct {
	Type  string `json:"type"`
	Cards []Card `json:"cards"`
}

func NewMeldPlay(cards []Card) MeldPlay {
	return MeldPlay{
		Type:  "PLAY_MELD",
		Cards: cards,
	}
}
func (m MeldPlay) GetName() AvailablePlay {
	return PLAY_MELD
}

func (m MeldPlay) GetCards() []Card {
	return m.Cards
}

type DrawCardPlay struct {
	Type string `json:"type"`
}

func NewDrawCardPlay() DrawCardPlay {
	return DrawCardPlay{Type: "DRAW_CARD"}
}

func (d DrawCardPlay) GetName() AvailablePlay {
	return DRAW_CARD
}
func (p DrawCardPlay) GetCards() []Card {
	log.Println("No cards used for this play")
	return nil
}

func IsValid(turnState *TurnState, play Play, outputProvider OutputProvider) bool {
	switch play.GetName() {

	case PLAY_MELD:
		return true // Meld is validated at client side

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
		log.Print("player :: !> Playing meld")
		for _, card := range play.GetCards() {
			if player.Hand.Contains(card) {
				player.Hand.RemoveCard(card)
				table.AddCard(&card)
			}
		}
		return

	case DRAW_CARD:
		log.Print("player :: !> Playing draw card")
		card := deck.DrawCard()
		player.Hand.AddCard(card)
		return

	case END_TURN:
		log.Print("player :: !> Passing turn")
		outputProvider.Write("message", "Passing turn")
		return

	case QUIT:
		return

	default:
		return
	}

}
