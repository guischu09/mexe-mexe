package server

import (
	"mexemexe/internal/engine"
	"mexemexe/internal/service"
	"net"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// SERVER_CAPACITY defines the maximum number of clients that can connect to the server
const SERVER_CAPACITY = 30

// Upgrader defines the websocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// GenerateUniqueID generates a unique ID
func GenerateUniqueID() string {
	return uuid.New().String()
}

// Server defines the game server struct
type Server struct {
	Clients  map[string]*Client
	Rooms    map[string]*GameRoom
	Capacity int
	mu       sync.Mutex
	config   *ServerConfig
	uuid     string
	logger   *service.GameLogger
}

// Server constructor
func NewServer(serverConfig *ServerConfig) *Server {
	uuid := GenerateUniqueID()
	return &Server{
		Clients:  make(map[string]*Client),
		Rooms:    make(map[string]*GameRoom, SERVER_CAPACITY),
		Capacity: SERVER_CAPACITY,
		config:   serverConfig,
		uuid:     uuid,
		logger:   service.NewLogger(serverConfig.logLevel, uuid),
	}
}

// AddClient adds a new client to the server
func (s *Server) AddClient(newClient *Client) {
	s.logger.Debugf("Adding client %s to server", newClient.Conn.RemoteAddr())
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Clients[newClient.UUID] = newClient
}

// RemoveClient removes a client from the server
func (s *Server) RemoveClient(uuid string) {
	s.logger.Debugf("Removing client with uuid %s from server", uuid)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Clients, uuid)
}

// AddRoom adds a new room to the server
func (s *Server) AddRoom(room *GameRoom) {
	s.logger.Debugf("Adding room %s to server", room.UUID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Rooms[room.UUID] = room
}

// RemoveClientFromRoom removes a client from a room
func (s *Server) RemoveClientFromRoom(room *GameRoom, client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Rooms[room.UUID].RemoveClient(client)
}

// RemoveRoom removes a room from the server
func (s *Server) RemoveRoom(room *GameRoom) {
	s.logger.Debugf("Removing room %s from server", room.UUID)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Rooms, room.UUID)
}

// SearchAvailableGameRoom searches for a game room with the specified number of players
func (s *Server) SearchAvailableGameRoom(numPlayers uint8) (*GameRoom, error) {
	s.mu.Lock()
	s.logger.Debugf("Searching among %d rooms", len(s.Rooms))

	candidateRooms := make([]*GameRoom, 0, len(s.Rooms))

	for id, room := range s.Rooms {
		s.logger.Debugf("Examining room %s", id)

		if room == nil {
			continue
		}

		room.mu.Lock()
		clientCount := len(room.Clients)
		gameStarted := room.GameStarted

		// Check if client list is nil
		if room.Clients == nil {
			s.logger.Debugf("Room %s has nil Clients slice!", id)
		}

		// Check individual clients
		for i, client := range room.Clients {
			if client == nil {
				s.logger.Debugf("Room %s has nil client at index %d!", id, i)
			} else {
				s.logger.Debugf("Room %s has nil client at index %d!", id, i)

			}
		}
		room.mu.Unlock()

		s.logger.Debugf("Room %s has %d clients, started: %v", id, clientCount, gameStarted)

		candidateRooms = append(candidateRooms, room)
	}
	s.mu.Unlock()

	var bestRoom *GameRoom
	var bestPlayerCount int = -1

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
		s.logger.Debugf("Found available game room: %s with %d/%d players",
			bestRoom.UUID, bestPlayerCount, numPlayers)
		return bestRoom, nil
	}
	return nil, nil
}

// GetCurrentCapacity returns the current number of connected clients
func (s *Server) GetCurrentCapacity() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Clients)
}

// IsAtMaximumCapacity returns true if the server is at maximum capacity
func (s *Server) IsAtMaximumCapacity() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Clients) >= s.Capacity
}

// AuthenticateUser checks if the user is authenticated
func (s *Server) AuthenticateUser(username string, ws *websocket.Conn) bool {
	return true
}

// parseRemoteAddr parses the remote address of a websocket connection into an IP and port
func parseRemoteAddr(addr string) (string, string) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr, ""
	}
	return host, port
}

// HandleConnections handles incoming websocket connections to the server
func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Establish a websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("error upgrading connection: %v", err)
		return
	}
	defer ws.Close()

	s.logger.Info("New connection established. Client connecting Address: " + r.RemoteAddr)

	// Read join message from client
	var joinMsg JoinServerMessage
	err = ws.ReadJSON(&joinMsg)
	if err != nil {
		s.logger.Errorf("error reading join message: %v", err)
		return
	}

	// Check if server is at maximum capacity
	if s.IsAtMaximumCapacity() {
		maxMsg := MaxCapacityMessage{
			Message: "Server is at maximum capacity. Please try again later.",
		}
		err = ws.WriteJSON(maxMsg)
		if err != nil {
			s.logger.Errorf("error writing max capacity message: %v", err)
		}
		return
	}

	// Authenticate user
	s.logger.Infof("Authenticating client from %s", r.RemoteAddr)
	if !s.AuthenticateUser(joinMsg.Username, ws) {
		s.logger.Infof("Authentication failed for client from %s", r.RemoteAddr)
		errorMsg := ErrorMessage{
			Message: "Authentication failed. Please create an account and try again.",
		}
		err = ws.WriteJSON(errorMsg)
		if err != nil {
			s.logger.Errorf("error writing auth error message: %v", err)
		}
		return
	}

	// Create a new client and register it in the server
	ip, port := parseRemoteAddr(r.RemoteAddr)
	uuid := GenerateUniqueID()
	newClient := NewClient(ip, port, joinMsg.Username, uuid, ws)
	s.AddClient(newClient)

	// Ensure client cleanup on function exit
	defer func() {
		s.RemoveClient(newClient.UUID)
		s.logger.Infof("Client %s disconnected and removed", newClient.UUID)
	}()

	// Send welcome message to client
	welcomeMsg := WelcomeMessage{
		Message:    "Welcome to mexe-mexe.com!",
		PlayerUUID: uuid,
	}
	err = ws.WriteJSON(welcomeMsg)
	if err != nil {
		s.logger.Errorf("error writing welcome message: %v", err)
		return
	}

	// Wait for start game message from client
	var startMsg StartGameMessage
	err = ws.ReadJSON(&startMsg)
	if err != nil {
		s.logger.Errorf("error reading start message: %v", err)
		return
	}

	// Handle the start message
	switch startMsg.Action {
	case "start":
		room, err := s.handleStartGame(newClient, ws)
		if err != nil {
			s.logger.Errorf("error handling start game: %v", err)
			return
		}
		if room != nil {
			defer func() {
				s.RemoveRoom(room)
				s.logger.Infof("Room %s was removed", room.UUID)
			}()
		}

	case "rejoin":
		s.handleRejoin()

	default:
		s.logger.Errorf("unknown action: %s", startMsg.Action)
		errorMsg := ErrorMessage{
			Message: "Unknown action. Please try again.",
		}
		err = ws.WriteJSON(errorMsg)
		if err != nil {
			s.logger.Errorf("error writing unknown action error: %v", err)
		}
		return
	}

	// After handling start game, just wait for the connection to close
	s.logger.Infof("Client %s setup complete, waiting for natural disconnection", newClient.UUID)
	select {}
}

// handleStartGame processes the start game request
func (s *Server) handleStartGame(client *Client, ws *websocket.Conn) (*GameRoom, error) {
	waitingMsg := JoinedGameRoomMessage{
		Message: "Searching for an available game room. Please wait ...",
	}
	err := ws.WriteJSON(waitingMsg)
	if err != nil {
		return nil, err
	}

	s.logger.Infof("Searching for an available game room to place client %s", client.UUID)
	room, err := s.SearchAvailableGameRoom(engine.NUM_PLAYERS)
	if err != nil {
		errorMsg := "Error finding game room: " + err.Error()
		return nil, ws.WriteJSON(errorMsg)
	}

	s.logger.Debugf("Found room: %v", room)

	// If room is available, add client to it
	if room != nil {
		return room, s.joinExistingRoom(client, room, ws)
	}

	// If no room is available, create a new one
	room = s.createNewRoom()
	s.AddRoom(room)
	s.logger.Debugf("Adding client %s to room %s", client.UUID, room.UUID)
	room.AddClient(client)
	s.logger.Debugf("Room has %d clients", len(room.Clients))
	s.logger.Debugf("Room clients usernames: %v", room.GetClientsUsername())

	joinedMsg := JoinedGameRoomMessage{
		Message: "Joined game room. Waiting for an opponent to join ...",
	}
	s.logger.Infof("Joined game room: %s. Waiting for an opponent to join ...", room.UUID)
	err = ws.WriteJSON(joinedMsg)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// joinExistingRoom adds client to an existing room
func (s *Server) joinExistingRoom(client *Client, room *GameRoom, ws *websocket.Conn) error {
	room.AddClient(client)

	joinedMsg := JoinedGameRoomMessage{
		Message: "Joined game room. Waiting for an opponent to join ...",
	}
	s.logger.Infof("Client %s joined game room: %s.", client.UUID, room.UUID)
	err := ws.WriteJSON(joinedMsg)
	if err != nil {
		return err
	}

	// If room is full, start the game
	if room.IsFull() {
		s.startGameInRoom(room)
	}

	return nil
}

// createNewRoom creates a new game room and adds the client
func (s *Server) createNewRoom() *GameRoom {
	s.logger.Debugf("No room available. Creating a new room.")
	room := NewGameRoom(s.config.logLevel)
	s.logger.Debugf("New room created with UUID: %s", room.UUID)
	return room
}

// startGameInRoom initializes and starts a game in the given room
func (s *Server) startGameInRoom(room *GameRoom) {
	playersUUIDs := room.GetClientsUUID()
	playersUsernames := room.GetClientsUsername()
	config := engine.NewGameConfig(playersUsernames, playersUUIDs)
	newGame := engine.NewGame(config, room.logger)
	room.AddGame(newGame)

	// Don't send separate "Game started!" message - the initial game state serves this purpose
	room.StartGame()
	s.logger.Infof("Game in room %s started!", room.UUID)
}

// handleRejoin processes rejoin requests (implement as needed)
func (s *Server) handleRejoin() {
	s.logger.Fatalf("Rejoin functionality not yet implemented.")
}

// Client defines a connected client
type Client struct {
	IP       string
	Port     string
	Conn     *websocket.Conn
	UUID     string
	Username string
}

// NewClient is Client constructor
func NewClient(ip string, port string, username string, uuid string, conn *websocket.Conn) *Client {
	client := Client{
		IP:       ip,
		Port:     port,
		Conn:     conn,
		UUID:     uuid,
		Username: username,
	}
	return &client
}
