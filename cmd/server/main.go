package main

import (
	"log"
	"mexemexe/pkg/engine"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type JoinServerMessage struct {
	Username string `json:"username"`
}

type GameStartMessage struct {
	Message string `json:"message"`
}

type Game struct {
}

type Room struct {
	ID   uint64
	Game GameRoom
}

type GameRoom struct {
	ID          uint64
	Players     []*Client
	Game        *engine.Game
	GameStarted bool
}
type Client struct {
	ID        uint64
	Username  string
	Connected bool
	Player    engine.Player
	Room      *GameRoom
	Conn      *websocket.Conn
}

func NewClient(username string, conn *websocket.Conn) *Client {
	return &Client{
		ID:        generateUniqueID(),
		Username:  username,
		Connected: true,
		Conn:      conn,
	}
}

func (c *Client) JoinRoom(room *GameRoom) bool {
	if len(room.Players) >= 2 {
		return false
	}

	// Add client to room
	room.Players = append(room.Players, c)
	c.Room = room

	// Create player in game engine
	c.Player = room.Game.AddPlayer(c.Username)

	// Check if game can start now
	if len(room.Players) == 2 && !room.GameStarted {
		room.StartGame()
	}

	return true
}

func (room *GameRoom) StartGame() {
	room.GameStarted = true
	room.Game.Start()

	// Notify players
	for _, client := range room.Players {
		// Send game start message through WebSocket
		gameStartMsg := struct {
			Type string `json:"type"`
			Data string `json:"data"`
		}{
			Type: "gameStart",
			Data: "The game has started!",
		}
		client.Conn.WriteJSON(gameStartMsg)
	}
}

func (c *Client) SetGameRoom(roomNumber uint64) {
}

func (c *Client) GetGameRoomStatus(gameRoom GameRoom) {
}

var clients = make(map[*websocket.Conn]*Client)
var gameRooms = make(map[uint64]*GameRoom)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	for {
		var msg JoinServerMessage
		err := ws.ReadJSON(&msg)

		clients[ws] = NewClient(msg.Username)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
	}

}

func main() {

	http.HandleFunc("/ws", handleConnections)

}
