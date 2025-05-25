package server

import (
	"mexemexe/internal/engine"
	"mexemexe/internal/service"
	"sync"
)

type GameRoom struct {
	UUID        string
	Game        *engine.Game
	Clients     []*Client
	NumPlayers  uint8
	GameStarted bool
	RoomChannel chan string
	mu          sync.Mutex
	logger      *service.GameLogger
}

func NewGameRoom(debugLevel int) *GameRoom {
	uuid := GenerateUniqueID()
	logger := service.NewLogger(debugLevel, uuid)
	gameRoom := GameRoom{
		UUID:        uuid,
		Game:        nil,
		Clients:     []*Client{},
		NumPlayers:  0,
		GameStarted: false,
		RoomChannel: make(chan string),
		logger:      logger,
	}
	logger.Debugf("New game room created with UUID: %s", uuid)
	return &gameRoom
}

func (g *GameRoom) AddGame(game *engine.Game) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Game = game
}

func (g *GameRoom) AddClient(Client *Client) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.logger.Debugf("Room %s before adding: %d clients", g.UUID, len(g.Clients))
	if !g.isFullLocked() {
		g.Clients = append(g.Clients, Client)
		g.NumPlayers = uint8(len(g.Clients))
		g.logger.Debugf("Room %s after adding: %d clients, client %s added",
			g.UUID, len(g.Clients), Client.UUID)
	} else {
		g.logger.Debugf("Room %s is full, can't add client %s",
			g.UUID, Client.UUID)
	}
}

func (g *GameRoom) RemoveClient(Client *Client) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i, client := range g.Clients {
		if client.UUID == Client.UUID {
			g.Clients = append(g.Clients[:i], g.Clients[i+1:]...)
			g.NumPlayers = uint8(len(g.Clients))
			break
		}
	}
}

func (g *GameRoom) GetClientsUsername() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	usernames := make([]string, len(g.Clients))
	for i, client := range g.Clients {
		usernames[i] = client.Username
	}
	return usernames
}

func (g *GameRoom) IsFull() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.NumPlayers > 2 {
		// At some point just log this error on the server
		g.logger.Fatalf("ERROR: Game room is full. Cannot add more clients. Num players: %d", g.NumPlayers)
	}
	return g.NumPlayers == 2
}

func (g *GameRoom) isFullLocked() bool {
	if g.NumPlayers > 2 {
		g.logger.Errorf("Game room is full. Cannot add more clients. Num players: %d", g.NumPlayers)
		return true
	}
	return g.NumPlayers == 2
}

func (g *GameRoom) StartGame() {
	g.logger.Infof("Game on room %s started!\n", g.UUID)
	g.mu.Lock()
	defer g.mu.Unlock()
	g.GameStarted = true

	var GameState GameStateMessage

	inputProvider := make([]engine.InputProvider, len(g.Clients))
	outputProvider := make([]engine.OutputProvider, len(g.Clients))

	firstPlayer := g.Clients[0]

	for i, client := range g.Clients {
		inputProvider[i] = NewWebsocketInputProvider(client.Conn, g.logger)
		outputProvider[i] = NewWebsocketOutputProvider(client.Conn, *g.logger)

		GameState = GameStateMessage{
			Table: g.Game.Table,
			Hand:  *g.Game.Players[i].Hand,
			Turn:  *engine.NewTurnState(firstPlayer.UUID),
		}
		err := client.Conn.WriteJSON(GameState)
		if err != nil {
			g.logger.Errorf("error writing to websocket: %v", err)
			return
		}
		g.logger.Debugf("Sent game state to client %s", client.UUID)
	}

	go g.Game.Start(inputProvider, outputProvider)
}
