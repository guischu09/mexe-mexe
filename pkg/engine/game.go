package engine

// Mexe-mexe rules:
// Wins a round the player that has 0 cards on their hands first.
// If the decks ends, the winner is the one with less points/cards on their hands.
// The game is played in turns.
// When a turn starts, one player can play a meld (sequence or book) on the table. If there are no melds available, the player must draw a card.
// If no meld is available, and no mexe-mexe is possible, the turn ends, and it is time for the next player to play.
// When a meld is played on the table (by you or your opponent), Mexe-mexe becomes possible.
// Mexe-mexe allows you to play with the cards (melds) on the table as if they were your own cards.
// And as long as the changes you do on the table are valid (melds are of a sequence or book), you can play them as if they were your own cards.

// Mecanicas:

// Draw a card = d
// Selecionar uma meld = m + #carta

// Jogar meld selecionada = P (ação ira printar as cartas)
// Jogar meld confirmar

// sugestao Mi: m1 m2 m3, m11
// sugestao do Gui: segurar m +  #carta. Para selecionar carta numero 11, pressionar 1 duas vezes rapidamente

// quando houver a mesa, usuario podera selecionar a mão "ctrl + h" ou a mesa "ctrl + t"

type GameConfig struct {
	Seed        uint64
	PlayersName []string
	NumPlayers  uint8
	NumCards    uint8
}

type Game struct {
	Config  GameConfig
	Deck    GameDeck
	Table   Table
	Players []Player
}

func NewGame(config GameConfig) Game {

	players := make([]Player, config.NumPlayers)
	deck := NewGameDeck(config.Seed)

	for i := 0; i < int(config.NumPlayers); i++ {
		newHand := NewHandFromDeck(&deck, config.NumCards)
		players[i] = NewPlayer(config.PlayersName[i], newHand, 0)
	}

	return Game{
		Config:  config,
		Deck:    deck,
		Players: players,
	}
}

func (g *Game) Start() {

	for c := 0; c > g.Deck.Size; c-- {

	}
}
