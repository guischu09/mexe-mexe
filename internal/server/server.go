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
	Clients  map[*websocket.Conn]*Client
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
		Clients:  make(map[*websocket.Conn]*Client),
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
	s.Clients[newClient.Conn] = newClient
}

// RemoveClient removes a client from the server
func (s *Server) RemoveClient(conn *websocket.Conn) {
	s.logger.Debugf("Removing client %s from server", conn.RemoteAddr())
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Clients, conn)
}

// AddRoom adds a new room to the server
func (s *Server) AddRoom(room *GameRoom) {
	s.logger.Debugf("Adding room %s to server", room.UUID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Rooms[room.UUID] = room
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

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Stablish a websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Errorf("error: %v", err)
	}
	defer ws.Close()

	s.logger.Info("New connection established. Client connecting Address: " + r.RemoteAddr)

	// Read join message from client
	var joinMsg JoinServerMessage
	err = ws.ReadJSON(&joinMsg)
	if err != nil {
		s.logger.Errorf("error: %v", err)
	}

	// Check if server is at maximum capacity
	if s.IsAtMaximumCapacity() {
		var maxMsg MaxCapacityMessage
		maxMsg.Message = "Server is at maximum capacity. Please try again later."
		err = ws.WriteJSON(maxMsg)
		if err != nil {
			s.logger.Errorf("error writing to websocket: %v", err)
			return
		}
		return
	}
	// Authenticate user:
	s.logger.Infof("Authenticating client from %s", r.RemoteAddr)
	if !s.AuthenticateUser(joinMsg.Username, ws) {
		s.logger.Infof("Authentication failed for client from %s", r.RemoteAddr)
		errorMsg := ErrorMessage{
			Message: "Authentication failed. Please create an account and try again.",
		}
		err = ws.WriteJSON(errorMsg)
		if err != nil {
			s.logger.Errorf("error writing to websocket: %v", err)
			return
		}
		return
	}

	// Create a new client and register it in the server
	ip, port := parseRemoteAddr(r.RemoteAddr)
	uuid := GenerateUniqueID()
	newClient := NewClient(ip, port, joinMsg.Username, uuid, ws)
	s.AddClient(newClient)

	// Send welcome message to client
	var welcomeMsg WelcomeMessage
	welcomeMsg.Message = "Welcome to mexe-mexe.com!"
	welcomeMsg.PlayerUUID = uuid
	err = ws.WriteJSON(welcomeMsg)
	if err != nil {
		s.logger.Errorf("error writing to websocket: %v", err)
		return
	}

	// Event loop
	for {
		// Read start game message from client
		var startMsg StartGameMessage
		err = ws.ReadJSON(&startMsg)
		if err != nil {
			s.logger.Errorf("error: %v", err)
			s.RemoveClient(ws)
			break
		}
		switch startMsg.Action {

		case "start":
			var waitingRoomMessage JoinedGameRoomMessage
			waitingRoomMessage.Message = "Searching for an available game room. Please wait ..."
			err = ws.WriteJSON(waitingRoomMessage)
			if err != nil {
				s.logger.Errorf("error writing to websocket: %v", err)
				s.RemoveClient(ws)
				return
			}
			s.logger.Infof("Searching for an available game room to place client %s\n", newClient.UUID)
			room, err := s.SearchAvailableGameRoom(engine.NUM_PLAYERS)
			s.logger.Debugf("Found room: %v\n", room)
			if err != nil {
				err = ws.WriteJSON("Error finding game room: " + err.Error())
				if err != nil {
					s.logger.Errorf("error writing to websocket: %v", err)
					return
				}
				s.RemoveClient(ws)
				continue
			}

			// If room is available (exists and not full) add client to room
			if room != nil {
				room.AddClient(newClient)
				var joinedRoomMsg JoinedGameRoomMessage
				joinedRoomMsg.Message = "Joined game room. Waiting for an opponent to join ..."
				s.logger.Infof("Joined game room: %s. Waiting for an opponent to join ...", room.UUID)
				err = ws.WriteJSON(joinedRoomMsg)
				if err != nil {
					s.logger.Errorf("error writing to websocket: %v", err)
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
					var gameStartedMsg GameStartedMessage
					gameStartedMsg.Message = "Game started!"
					err = ws.WriteJSON(gameStartedMsg)
					if err != nil {
						s.logger.Errorf("error writing to websocket: %v", err)
						s.RemoveClient(ws)
						return
					}

				}

			}

			// If no room is available, create a new one
			if room == nil {
				s.logger.Debugf("No room available. Creating a new room.")
				room = NewGameRoom(s.config.logLevel)
				s.logger.Debugf("New room created with UUID: %s", room.UUID)
				s.AddRoom(room)
				s.logger.Debugf("Room added to server")

				// Check that room exists in the map
				s.mu.Lock()
				if existingRoom, ok := s.Rooms[room.UUID]; ok {
					s.logger.Debugf("Room %s found in map after adding", room.UUID)
					if existingRoom != room {
						s.logger.Warningf("Room pointers don't match!")
					}
				} else {
					s.logger.Warningf("Room %s not found in map after adding!", room.UUID)
				}
				s.mu.Unlock()

				s.logger.Debugf("About to add client %s to room %s", newClient.UUID, room.UUID)
				room.AddClient(newClient)
				s.logger.Debugf("After adding client, room has %d clients", len(room.Clients))

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
					s.logger.Debugf("Client %s not found in room after adding!", newClient.UUID)
				} else {
					s.logger.Debugf("Client %s confirmed in room", newClient.UUID)
				}

				var joinedRoomMsg JoinedGameRoomMessage
				joinedRoomMsg.Message = "Joined game room. Waiting for an opponent to join ..."
				s.logger.Infof("Joined game room: %s. Waiting for an opponent to join ...", room.UUID)
				err = ws.WriteJSON(joinedRoomMsg)
				if err != nil {
					s.logger.Errorf("error writing to websocket: %v", err)
					return
				}
			}
		case "rejoin":

		}

	}

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
