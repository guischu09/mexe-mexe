package main

import (
	"errors"
	"fmt"
	"log"
	"mexemexe/internal/engine"
	"mexemexe/internal/server"
	"net"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const SERVER_CAPACITY = 30

type Server struct {
	Clients map[*websocket.Conn]*Client
	Rooms   map[string]*GameRoom
	mu      sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Clients: make(map[*websocket.Conn]*Client),
		Rooms:   make(map[string]*GameRoom, SERVER_CAPACITY),
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
	if len(s.Rooms) >= SERVER_CAPACITY {
		return nil, errors.New("server at maximum capacity, cannot create new room")
	}
	for _, room := range s.Rooms {
		if room != nil && !room.GameStarted && len(room.Clients) < int(numPlayers) {
			log.Printf("Found available game room: %s", room.UUID)
			return room, nil
		}
	}
	return nil, nil
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
	return &Client{
		IP:       ip,
		Port:     port,
		Conn:     conn,
		UUID:     uuid,
		Username: username,
	}
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
	uuid := generateUniqueID()
	return &GameRoom{
		UUID:        uuid,
		Game:        nil,
		Clients:     []*Client{},
		NumPlayers:  0,
		GameStarted: false,
		RoomChannel: make(chan string),
	}
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
		log.Fatalf("ERROR: Game room is full. Cannot add more clients. Num players: %d", g.NumPlayers)
	}
	return g.NumPlayers == 2
}

func (g *GameRoom) StartGame() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.GameStarted = true
	g.RoomChannel <- "Game started!"
	go g.Game.Start()
}

// TODO: Implement
func generateUniqueID() string {
	return uuid.New().String()
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
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer ws.Close()

	log.Println("New connection established. Connecting Address: " + r.RemoteAddr)

	var msg server.JoinServerMessage
	err = ws.ReadJSON(&msg)
	if err != nil {
		log.Printf("error: %v", err)
	}

	ip, port := parseRemoteAddr(r.RemoteAddr)
	uuid := generateUniqueID()

	newClient := NewClient(ip, port, msg.Username, uuid, ws)

	s.AddClient(ws, ip, port, msg.Username, uuid)
	err = ws.WriteJSON("Welcome to mexe-mexe.com!")
	if err != nil {
		log.Printf("error writing to websocket: %v", err)
		return
	}

	// Event loop
	for {
		var msg server.StartGameMessage
		err = ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			s.RemoveClient(ws)
			break
		}
		switch msg.Type {

		case "rejoin":

		case "start":
			err = ws.WriteJSON("Searching for an available game room. Please wait ...")
			if err != nil {
				log.Printf("error writing to websocket: %v", err)
				return
			}
			fmt.Printf("Searching for an available game room to place client %s\n", newClient.UUID)
			room, err := s.SearchAvailableGameRoom(engine.NUM_PLAYERS)
			if err != nil {
				err = ws.WriteJSON("Error finding game room: " + err.Error())
				if err != nil {
					log.Printf("error writing to websocket: %v", err)
					return
				}
				continue
			}

			// If room is available and not full add client to room
			if room != nil {
				fmt.Println("Found available room with UUID: " + room.UUID)
				room.AddClient(newClient)
				msg := fmt.Sprintf("Joined game room: %s. Waiting for an opponent to join ...", room.UUID)
				err = ws.WriteJSON(msg)
				if err != nil {
					log.Printf("error writing to websocket: %v", err)
					return
				}
				// If room is full, start game
				if room.IsFull() {
					playersUsernames := room.GetClientsUsername()
					config := engine.NewGameConfig(playersUsernames)
					newGame := engine.NewGame(config)
					room.AddGame(newGame)
					room.StartGame()
					// broadcastToRoom(room, "Game started!", "System")
				}

				if len(room.Clients) >= 2 {
					room.GameStarted = true
				}

			}

			// If no room is available, create a new one
			if room == nil {
				fmt.Println("No room available. Creating a new room.")
				room = NewGameRoom()
				s.AddRoom(room)
				room.AddClient(newClient)
				msg := fmt.Sprintf("Joined game room: %s. Waiting for an opponent to join ...", room.UUID)
				err = ws.WriteJSON(msg)
				if err != nil {
					log.Printf("error writing to websocket: %v", err)
					return
				}
			}

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
