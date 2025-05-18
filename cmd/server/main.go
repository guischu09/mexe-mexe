package main

import (
	"log"
	"mexemexe/internal/engine"
	"mexemexe/internal/server"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const SERVER_CAPACITY = 30

type Server struct {
	Clients  map[*websocket.Conn]*Client
	Rooms    map[string]*GameRoom
	Capacity int
	mu       sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Clients:  make(map[*websocket.Conn]*Client),
		Rooms:    make(map[string]*GameRoom, SERVER_CAPACITY),
		Capacity: SERVER_CAPACITY,
	}
}

func (s *Server) AddClient(newClient *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Clients[newClient.Conn] = newClient
}

func (s *Server) RemoveClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Clients, conn)
}

func (s *Server) AddRoom(room *GameRoom) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Rooms[room.UUID] = room
}

func (s *Server) RemoveRoom(room *GameRoom) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Rooms, room.UUID)
}

func (s *Server) SearchAvailableGameRoom(numPlayers uint8) (*GameRoom, error) {
	s.mu.Lock()
	log.Printf("DEBUG: SearchRoom - Searching among %d rooms", len(s.Rooms))

	candidateRooms := make([]*GameRoom, 0, len(s.Rooms))
	for id, room := range s.Rooms {
		log.Printf("DEBUG: SearchRoom - Examining room %s", id)

		// Check if room is nil
		if room == nil {
			log.Printf("DEBUG: SearchRoom - Room %s is nil!", id)
			continue
		}

		// Look at room details
		room.mu.Lock()
		clientCount := len(room.Clients)
		gameStarted := room.GameStarted

		// Check if client list is nil
		if room.Clients == nil {
			log.Printf("DEBUG: SearchRoom - Room %s has nil Clients slice!", id)
		}

		// Check individual clients
		for i, client := range room.Clients {
			if client == nil {
				log.Printf("DEBUG: SearchRoom - Room %s has nil client at index %d!", id, i)
			} else {
				log.Printf("DEBUG: SearchRoom - Room %s has client %s at index %d",
					id, client.UUID, i)
			}
		}

		room.mu.Unlock()

		log.Printf("DEBUG: SearchRoom - Room %s has %d clients, started: %v",
			id, clientCount, gameStarted)

		candidateRooms = append(candidateRooms, room)
	}
	s.mu.Unlock()

	var bestRoom *GameRoom
	var bestPlayerCount int = -1

	// Now examine each room without holding the server lock
	for _, room := range candidateRooms {
		room.mu.Lock()
		clientCount := len(room.Clients)
		isAvailable := !room.GameStarted &&
			clientCount < int(numPlayers) &&
			clientCount > 0
		room.mu.Unlock()

		if isAvailable && clientCount > bestPlayerCount {
			bestRoom = room
			bestPlayerCount = clientCount
		}
	}

	if bestRoom != nil {
		log.Printf("Found available game room: %s with %d/%d players",
			bestRoom.UUID, bestPlayerCount, numPlayers)
		return bestRoom, nil
	}

	return nil, nil
}

func (s *Server) GetCurrentCapacity() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Clients)
}

func (s *Server) IsAtMaximumCapacity() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Clients) >= s.Capacity
}

func (s *Server) AuthenticateUser(username string, ws *websocket.Conn) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	IP       string
	Port     string
	Conn     *websocket.Conn
	UUID     string
	Username string
}

func NewClient(ip string, port string, username string, uuid string, conn *websocket.Conn) *Client {
	client := Client{
		IP:       ip,
		Port:     port,
		Conn:     conn,
		UUID:     uuid,
		Username: username,
	}
	log.Println("New client created with UUID: " + uuid)
	return &client
}

type GameRoom struct {
	UUID        string
	Game        *engine.Game
	Clients     []*Client
	NumPlayers  uint8
	GameStarted bool
	RoomChannel chan string
	mu          sync.Mutex
}

func NewGameRoom() *GameRoom {
	uuid := server.GenerateUniqueID()
	gameRoom := GameRoom{
		UUID:        uuid,
		Game:        nil,
		Clients:     []*Client{},
		NumPlayers:  0,
		GameStarted: false,
		RoomChannel: make(chan string),
	}
	log.Println("New game room created with UUID: " + uuid)
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
	log.Printf("DEBUG: AddClient - Room %s before adding: %d clients", g.UUID, len(g.Clients))
	if !g.isFullLocked() {
		g.Clients = append(g.Clients, Client)
		g.NumPlayers = uint8(len(g.Clients))
		log.Printf("DEBUG: AddClient - Room %s after adding: %d clients, client %s added",
			g.UUID, len(g.Clients), Client.UUID)
	} else {
		log.Printf("DEBUG: AddClient - Room %s is full, can't add client %s",
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
		log.Fatalf("ERROR: Game room is full. Cannot add more clients. Num players: %d", g.NumPlayers)
	}
	return g.NumPlayers == 2
}

func (g *GameRoom) isFullLocked() bool {
	if g.NumPlayers > 2 {
		log.Printf("ERROR: Game room is full. Cannot add more clients. Num players: %d", g.NumPlayers)
		return true
	}
	return g.NumPlayers == 2
}

func (g *GameRoom) StartGame() {
	log.Printf("Game on room %s started!\n", g.UUID)
	g.mu.Lock()
	defer g.mu.Unlock()
	g.GameStarted = true
	// g.RoomChannel <- "Game started!"

	inputProvider := make([]engine.InputProvider, len(g.Clients))
	outputProvider := make([]engine.OutputProvider, len(g.Clients))
	for i, client := range g.Clients { // i is index, client is the value
		inputProvider[i] = NewWebsocketInputProvider(client.Conn)
		outputProvider[i] = NewWebsocketOutputProvider(client.Conn)
	}
	go g.Game.Start(inputProvider, outputProvider)
}

type WebsocketOutputProvider struct {
	conn *websocket.Conn
}

func NewWebsocketOutputProvider(conn *websocket.Conn) WebsocketOutputProvider {
	return WebsocketOutputProvider{
		conn: conn,
	}
}

func (w WebsocketOutputProvider) Write(messageType string, data interface{}) {
	log.Printf("DEBUG: Write - Writing message type %s", messageType)
}

type WebsocketInputProvider struct {
	conn *websocket.Conn
}

func NewWebsocketInputProvider(conn *websocket.Conn) WebsocketInputProvider {
	return WebsocketInputProvider{
		conn: conn,
	}
}

func (w WebsocketInputProvider) GetPlay(table *engine.Table, hand *engine.Hand, playerName string, turnState *engine.TurnState) engine.Play {
	// Create a properly populated game state message
	gameStateMsg := server.GameStateMessage{
		Table: *table,
		Hand:  *hand,
		Turn:  *turnState,
	}

	// Send the game state to the client
	err := w.conn.WriteJSON(&gameStateMsg)
	if err != nil {
		log.Printf("error writing to websocket: %v", err)
		log.Println("Closing connection")
		w.conn.Close()
		// Return a "quit" play instead of nil
		return engine.NewQuitPlay("connection_lost")
	}

	// Read the client's move
	var gamePlayMsg server.GamePlayMessage
	err = w.conn.ReadJSON(&gamePlayMsg)
	if err != nil {
		log.Printf("error reading from websocket: %v", err)
		log.Println("Closing connection")
		w.conn.Close()
		// Return a "quit" play instead of nil
		return engine.NewQuitPlay("connection_lost")
	}

	// Handle case where Play might be nil
	if gamePlayMsg.Play == nil {
		log.Printf("warning: received nil play from client, treating as quit")
		return engine.NewQuitPlay("nil_play")
	}

	return gamePlayMsg.Play
}

// parseRemoteAddr parses the remote address of a websocket connection into an IP and port
func parseRemoteAddr(addr string) (string, string) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr, ""
	}
	return host, port
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Stablish a websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer ws.Close()

	log.Println("New connection established. Client connecting Address: " + r.RemoteAddr)

	// Read join message from client
	var joinMsg server.JoinServerMessage
	err = ws.ReadJSON(&joinMsg)
	if err != nil {
		log.Printf("error: %v", err)
	}

	// Check if server is at maximum capacity
	if s.IsAtMaximumCapacity() {
		var maxMsg server.MaxCapacityMessage
		maxMsg.Message = "Server is at maximum capacity. Please try again later."
		err = ws.WriteJSON(maxMsg)
		if err != nil {
			log.Printf("error writing to websocket: %v", err)
			return
		}
		return
	}
	// Authenticate user:
	log.Printf("Authenticating client from %s", r.RemoteAddr)
	if !s.AuthenticateUser(joinMsg.Username, ws) {
		log.Printf("Authentication failed for client from %s", r.RemoteAddr)
		errorMsg := server.ErrorMessage{
			Message: "Authentication failed. Please create an account and try again.",
		}
		err = ws.WriteJSON(errorMsg)
		if err != nil {
			log.Printf("error writing to websocket: %v", err)
			return
		}
		return
	}

	// Create a new client and register it in the server
	ip, port := parseRemoteAddr(r.RemoteAddr)
	uuid := server.GenerateUniqueID()
	newClient := NewClient(ip, port, joinMsg.Username, uuid, ws)
	s.AddClient(newClient)

	// Send welcome message to client
	var welcomeMsg server.WelcomeMessage
	welcomeMsg.Message = "Welcome to mexe-mexe.com!"
	welcomeMsg.PlayerUUID = uuid
	err = ws.WriteJSON(welcomeMsg)
	if err != nil {
		log.Printf("error writing to websocket: %v", err)
		return
	}

	// Event loop
	for {
		// Read start game message from client
		var startMsg server.StartGameMessage
		err = ws.ReadJSON(&startMsg)
		if err != nil {
			log.Printf("error: %v", err)
			s.RemoveClient(ws)
			break
		}
		switch startMsg.Action {

		case "start":
			var waitingRoomMessage server.JoinedGameRoomMessage
			waitingRoomMessage.Message = "Searching for an available game room. Please wait ..."
			err = ws.WriteJSON(waitingRoomMessage)
			if err != nil {
				log.Printf("error writing to websocket: %v", err)
				s.RemoveClient(ws)
				return
			}
			log.Printf("Searching for an available game room to place client %s\n", newClient.UUID)
			room, err := s.SearchAvailableGameRoom(engine.NUM_PLAYERS)
			log.Printf("DEBUG: Room: %v\n", room)
			if err != nil {
				err = ws.WriteJSON("Error finding game room: " + err.Error())
				if err != nil {
					log.Printf("error writing to websocket: %v", err)
					return
				}
				s.RemoveClient(ws)
				continue
			}

			// If room is available (exists and not full) add client to room
			if room != nil {
				room.AddClient(newClient)
				var joinedRoomMsg server.JoinedGameRoomMessage
				joinedRoomMsg.Message = "Joined game room. Waiting for an opponent to join ..."
				log.Printf("Joined game room: %s. Waiting for an opponent to join ...", room.UUID)
				err = ws.WriteJSON(joinedRoomMsg)
				if err != nil {
					log.Printf("error writing to websocket: %v", err)
					s.RemoveClient(ws)
					return
				}
				// If room is full, start game
				if room.IsFull() {
					playersUsernames := room.GetClientsUsername()
					config := engine.NewGameConfig(playersUsernames)
					newGame := engine.NewGame(config)
					room.AddGame(newGame)
					room.StartGame()
					var gameStartedMsg server.GameStartedMessage
					gameStartedMsg.Message = "Game started!"
					err = ws.WriteJSON(gameStartedMsg)
					if err != nil {
						log.Printf("error writing to websocket: %v", err)
						s.RemoveClient(ws)
						return
					}

				}

			}

			// If no room is available, create a new one
			if room == nil {
				log.Println("No room available. Creating a new room.")
				room = NewGameRoom()
				log.Printf("DEBUG: New room created with UUID: %s", room.UUID)
				s.AddRoom(room)
				log.Printf("DEBUG: Room added to server")

				// Check that room exists in the map
				s.mu.Lock()
				if existingRoom, ok := s.Rooms[room.UUID]; ok {
					log.Printf("DEBUG: Room %s found in map after adding", room.UUID)
					if existingRoom != room {
						log.Printf("DEBUG: WARNING - Room pointers don't match!")
					}
				} else {
					log.Printf("DEBUG: WARNING - Room %s not found in map after adding!", room.UUID)
				}
				s.mu.Unlock()

				log.Printf("DEBUG: About to add client %s to room %s", newClient.UUID, room.UUID)
				room.AddClient(newClient)
				log.Printf("DEBUG: After adding client, room has %d clients", len(room.Clients))

				// Verify client was added
				room.mu.Lock()
				clientFound := false
				for _, c := range room.Clients {
					if c.UUID == newClient.UUID {
						clientFound = true
						break
					}
				}
				room.mu.Unlock()

				if !clientFound {
					log.Printf("DEBUG: WARNING - Client %s not found in room after adding!", newClient.UUID)
				} else {
					log.Printf("DEBUG: Client %s confirmed in room", newClient.UUID)
				}

				var joinedRoomMsg server.JoinedGameRoomMessage
				joinedRoomMsg.Message = "Joined game room. Waiting for an opponent to join ..."
				log.Printf("Joined game room: %s. Waiting for an opponent to join ...", room.UUID)
				err = ws.WriteJSON(joinedRoomMsg)
				if err != nil {
					log.Printf("error writing to websocket: %v", err)
					return
				}
			}
		case "rejoin":

		}

	}

}

func main() {

	server := NewServer()
	http.HandleFunc("/ws", server.handleConnections)

	log.Println("HTTP server started on :8888")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
