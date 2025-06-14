package engine

import (
	"fmt"
	"log"
	"math/rand/v2"
	"mexemexe/internal/service"
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
const NUM_PLAYERS = 2
const NUM_CARDS = 21

type GameOptions struct {
	NumPlayers uint8
	NumCards   uint8
}

type GameConfig struct {
	Seed              uint64
	PlayersName       []string
	PlayersUUID       []string
	NumPlayers        uint8
	NumCards          uint8
	RandomPlayerOrder bool
	TotalCards        uint8
}

func NewGameConfig(playersNames []string, playersUUID []string) *GameConfig {
	gameConfig := GameConfig{
		Seed:              UNIQUE_SHUFFLE_SEED,
		PlayersName:       playersNames,
		PlayersUUID:       playersUUID,
		NumPlayers:        uint8(len(playersNames)),
		NumCards:          NUM_CARDS,
		RandomPlayerOrder: true,
		TotalCards:        uint8(TOTAL_DECK_SIZE),
	}
	return &gameConfig
}

type Game struct {
	Config  *GameConfig
	Deck    *Deck
	Table   Table
	Players []Player
	logger  *service.GameLogger
}

func NewEmptyGame(config *GameConfig, logger *service.GameLogger) *Game {
	return &Game{
		Config:  config,
		Deck:    NewDeck(config.Seed),
		Table:   Table{},
		Players: nil,
		logger:  logger,
	}
}

func NewGame(config *GameConfig, logger *service.GameLogger) *Game {

	// Implement input validation for a game to start, to avoid starting games with 10 players.
	players := make([]Player, config.NumPlayers)
	deck := NewDeck(config.Seed)

	for i, uuid := range config.PlayersUUID {
		newHand := NewHandFromDeck(deck, config.NumCards)
		players[i] = NewPlayer(config.PlayersName[i], newHand, uuid, INITIAL_POINTS)
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
		logger:  logger,
	}
}

func (g *Game) AddPlayer(player Player) {
	g.Players = append(g.Players, player)
}

func (g *Game) ShufflePlayers() {
	pcgSource := rand.NewPCG(g.Deck.Seed, g.Deck.Seed)
	rng := rand.New(pcgSource)
	rng.Shuffle(len(g.Players), func(i, j int) {
		g.Players[i], g.Players[j] = g.Players[j], g.Players[i]
	})
}

func (g *Game) Start(inputProvider []InputProvider, outputProvider []OutputProvider) bool {

	g.logger.Infof("Game started!\r\n")
	g.logger.Infof("Players: %v\r\n", len(g.Players))
	g.logger.Infof("Deck: %v\r\n", g.Deck.Size)
	g.logger.Infof("Table: %v\r\n", g.Table.Size)

	for g.Deck.Size > 0 {
		for i := range g.Players {
			g.ValidadeGame()
			player := &g.Players[i]
			availablePlay := player.PlayTurn(g.Deck, &g.Table, inputProvider[i], outputProvider)

			switch availablePlay {
			case QUIT:
				g.logger.Infof("Player %s quits", player.Name)
				return false

			case END_TURN:
				g.logger.Infof("Player %s ends turn", player.Name)

			case DRAW_CARD:
				g.logger.Infof("Player %s drawed a card", player.Name)

			case PLAY_MELD:
				g.logger.Infof("Player %s played a meld", player.Name)
			}
			if player.Hand.Size == 0 {
				g.logger.Infof("Player %s wins!", player.Name)
				return true
			}
		}
	}
	// Broadcast to all players
	// outputProvider.Write("message", "Deck is empty! Game over!")
	fmt.Println("Deck is empty! Game over!")
	// g.ComputePoints()
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

func (g *Game) ValidadeGame() {
	numberCardsWithPlayers := 0
	for i := range g.Players {
		numberCardsWithPlayers += len(g.Players[i].Hand.Cards)
	}
	totalCardsGame := numberCardsWithPlayers + g.Deck.Size + len(g.Table.Cards)

	if uint8(totalCardsGame) == uint8(g.Config.TotalCards) {
		return
	}
	log.Fatalf("ERROR: Card leek. Current total cards: %d, expected: %d", totalCardsGame, g.Config.TotalCards)
}

// SendStateToPlayers sends the current state to all players via outputProviders
func SendStateToPlayers(outputProviders []OutputProvider, table Table, Hand Hand, turnState TurnState) {
	for _, outputProvider := range outputProviders {
		if turnState.PlayerUUID == outputProvider.GetUUID() {
			outputProvider.SendState(table, Hand, turnState)
		} else {
			outputProvider.SendState(table, EMPTY_HAND, turnState)
		}
	}
}
