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

func (s *Server) AddClient(conn *websocket.Conn, ip string, port string, username string, uuid string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Clients[conn] = NewClient(ip, port, username, uuid, conn)
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
	defer s.mu.Unlock()
	// if len(s.Rooms) >= SERVER_CAPACITY {
	// 	return nil, errors.New("server at maximum capacity, cannot create new room")
	// }
	for _, room := range s.Rooms {
		if room != nil && !room.GameStarted && len(room.Clients) < int(numPlayers) {
			log.Printf("Found available game room: %s", room.UUID)
			return room, nil
		}
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
	if !g.IsFull() {
		g.Clients = append(g.Clients, Client)
		g.NumPlayers = uint8(len(g.Clients))
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

func (g *GameRoom) StartGame() {
	log.Printf("Game on room %s started!\n", g.UUID)
	g.mu.Lock()
	defer g.mu.Unlock()
	g.GameStarted = true
	g.RoomChannel <- "Game started!"

	inputProvider := make([]engine.InputProvider, len(g.Clients))
	outputProvider := make([]engine.OutputProvider, len(g.Clients))
	for u := range g.Clients {
		inputProvider[u] = NewWebsocketInputProvider(g.Clients[u].Conn)
		outputProvider[u] = NewWebsocketOutputProvider(g.Clients[u].Conn)
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

	var gameStateMsg server.GameStateMessage
	err := w.conn.WriteJSON(&gameStateMsg)
	if err != nil {
		log.Printf("error writing to websocket: %v", err)
		log.Println("Closing connection")
		w.conn.Close()
		return nil
	}
	var gamePlayMsg server.GamePlayMessage
	err = w.conn.ReadJSON(&gamePlayMsg)
	if err != nil {
		log.Printf("error reading from websocket: %v", err)
		log.Println("Closing connection")
		w.conn.Close()
		return nil
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
	s.AddClient(ws, ip, port, joinMsg.Username, uuid)

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

				if len(room.Clients) >= 2 {
					room.GameStarted = true
				}

			}

			// If no room is available, create a new one
			if room == nil {
				log.Println("No room available. Creating a new room.")
				room = NewGameRoom()
				s.AddRoom(room)
				room.AddClient(newClient)
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
