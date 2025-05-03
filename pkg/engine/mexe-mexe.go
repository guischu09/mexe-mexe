package engine

import (
	"fmt"
	"math/rand/v2"
)

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

const INITIAL_POINTS uint32 = 0

type GameConfig struct {
	Seed              uint64
	PlayersName       []string
	NumPlayers        uint8
	NumCards          uint8
	RandomPlayerOrder bool
}

type Game struct {
	Config  GameConfig
	Deck    Deck
	Table   Table
	Players []Player
}

func NewGame(config GameConfig) *Game {

	// Implement input validation for a game to start, to avoid starting games with 10 players.
	players := make([]Player, config.NumPlayers)
	deck := NewDeck(config.Seed)

	for i := 0; i < int(config.NumPlayers); i++ {
		newHand := NewHandFromDeck(&deck, config.NumCards)
		players[i] = NewPlayer(config.PlayersName[i], newHand, INITIAL_POINTS)
	}

	if config.RandomPlayerOrder {
		pcgSource := rand.NewPCG(deck.Seed, deck.Seed)
		rng := rand.New(pcgSource)
		rng.Shuffle(len(players), func(i, j int) {
			players[i], players[j] = players[j], players[i]
		})
	}

	return &Game{
		Config:  config,
		Deck:    deck,
		Players: players,
	}
}

func (g *Game) Start() bool {

	fmt.Printf("Game started!\r\n")

	inputProvider := TerminalInputProvider{}
	outputProvider := TerminalOutputProvider{}

	for g.Deck.Size > 0 {
		for i := range g.Players {
			player := &g.Players[i]

			g.Print(player)

			availablePlay := player.PlayTurn(&g.Deck, &g.Table, &inputProvider, &outputProvider)

			switch availablePlay {

			case QUIT:
				outputProvider.Write("message", "Player "+player.Name+" quits")
				outputProvider.Write("message", "Game Over!")
				return false

			case END_TURN:
				outputProvider.Write("message", "Player "+player.Name+" ends turn")

			case DRAW_CARD:
				outputProvider.Write("message", "Player "+player.Name+" drawed a card")

			case PLAY_MELD:
				outputProvider.Write("message", "Player "+player.Name+" played a meld")
			}
			if player.Hand.Size == 0 {
				outputProvider.Write("message", "Player "+player.Name+" wins!")
				return true
			}
		}
	}
	outputProvider.Write("message", "Deck is empty! Game over!")
	g.ComputePoints()
	return true
}

func (g *Game) ComputePoints() {
	// TODO: Implement
}

func (g *Game) Close() {
	fmt.Printf("Game closed!\r\n")
}

func (g *Game) Print(player *Player) {
	fmt.Printf("%s's Hand:\r\n", player.Name)
	player.Hand.Print()
	fmt.Printf("Deck size: %d\r\n", g.Deck.Size)
}
