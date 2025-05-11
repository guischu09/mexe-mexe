package server

// import (
// 	"fmt"
// 	"log"
// 	"mexemexe/internal/engine"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// const CAPACITY = 30
// const MAX_PLAYERS = 2

// var clients = make(map[*websocket.Conn]*Client)
// var serverGameRooms = make([]*GameRoom, 0, CAPACITY)

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// type Client struct {
// 	IP       string
// 	Port     string
// 	Conn     *websocket.Conn
// 	Username string
// 	UUID     string
// }

// func NewClient(ip string, port string, username string, uuid string) *Client {
// 	return &Client{
// 		IP:       ip,
// 		Port:     port,
// 		UUID:     uuid,
// 		Username: username,
// 	}
// }

// type GameRoom struct {
// 	UUID        string
// 	Game        *engine.Game
// 	Players     []Client
// 	NumPlayers  uint8
// 	GameStarted bool
// 	RoomChannel chan string
// }

// func (g *GameRoom) IsFull() bool {
// 	return g.NumPlayers >= 2
// }

// func (g *GameRoom) StartGame() {
// 	g.GameStarted = true
// }

// func (g *GameRoom) HasSpace() bool {
// 	return g.NumPlayers < 2
// }

// // TODO: Implement
// func generateUniqueID() string {
// 	return "1234567890"
// }

// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("error: %v", err)
// 	}
// 	defer ws.Close()

// 	log.Println("New connection established. Connecting Address: " + r.RemoteAddr)

// 	var msg JoinServerMessage
// 	err = ws.ReadJSON(&msg)
// 	if err != nil {
// 		log.Printf("error: %v", err)
// 		delete(clients, ws)
// 	}

// 	newPlayer := NewClient(r.RemoteAddr, r.RemoteAddr, msg.Username, generateUniqueID())

// 	clients[ws] = newPlayer
// 	err = ws.WriteJSON("Welcome to mexe-mexe.com!")
// 	if err != nil {
// 		log.Printf("error writing to websocket: %v", err)
// 		return
// 	}

// 	// Event loop
// 	for {
// 		var msg GameMessage
// 		err = ws.ReadJSON(&msg)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			delete(clients, ws)
// 			close(newPlayer.ChatChannel)
// 			break
// 		}
// 		switch msg.Type {

// 		case "chat":

// 		case "start":
// 			err = ws.WriteJSON("Searching for an available game room ...")
// 			if err != nil {
// 				log.Printf("error writing to websocket: %v", err)
// 				return
// 			}
// 			Nplayers := uint8(2)
// 			room, err := searchAvailableGameRoom(Nplayers)
// 			if err != nil {
// 				err = ws.WriteJSON("Error finding game room: " + err.Error())
// 				if err != nil {
// 					log.Printf("error writing to websocket: %v", err)
// 					return
// 				}
// 				continue
// 			}

// 			room.AddPlayer(newPlayer)

// 			err = ws.WriteJSON("Joined game room: " + room.ID)
// 			if err != nil {
// 				log.Printf("error writing to websocket: %v", err)
// 				return
// 			}

// 			if room.IsFull() {
// 				room.StartGame()
// 				broadcastToRoom(room, "Game started!", "System")
// 			}

// 			if len(room.Players) >= 2 {
// 				room.GameStarted = true

// 			}

// 		case "quit":
// 			fmt.Println("Not implemented")
// 		}

// 	}

// }

// func broadcastToRoom(room *GameRoom, message string, sender string) {
// 	for _, player := range room.Players {
// 		player.ChatChannel <- sender + ": " + message
// 	}
// }

// func searchAvailableGameRoom(numPlayers uint8) (*GameRoom, error) {
// 	for _, room := range serverGameRooms {
// 		if room != nil && !room.GameStarted && len(room.Players) < int(numPlayers) {
// 			log.Printf("Found available game room: %s", room.ID)
// 			return room, nil
// 		}
// 	}

// 	if len(serverGameRooms) >= CAPACITY {
// 		return nil, errors.New("server at maximum capacity, cannot create new room")
// 	}

// 	log.Println("No available game rooms found. Creating a new game room.")
// 	newRoom := NewGameRoom()
// 	serverGameRooms = append(serverGameRooms, newRoom)
// 	return newRoom, nil
// }

// func main() {

// 	http.HandleFunc("/ws", handleConnections)

// 	log.Println("HTTP server started on :8888")
// 	err := http.ListenAndServe(":8888", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}

// }
