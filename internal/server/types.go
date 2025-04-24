package server

// import "github.com/gorilla/websocket"

// type JoinServerMessage struct {
// 	Username string `json:"username"`
// }

// type GameMessage struct {
// 	Type    string `json:"type"`
// 	Context string `json:"context"`
// }

// type MexeMexeGame struct {
// }

// type Room struct {
// 	ID   string
// 	Game GameRoom
// }

// type GameRoom struct {
// 	Game        *MexeMexeGame
// 	Players     []*Player
// 	ID          string
// 	NumPlayers  uint8
// 	GameStarted bool
// 	RoomChannel chan string
// }

// func (g *GameRoom) AddPlayer(player *Player) {
// 	g.Players = append(g.Players, player)
// 	g.NumPlayers = uint8(len(g.Players))
// 	player.Room = g
// }

// func (g *GameRoom) IsFull() bool {
// 	return g.NumPlayers >= 2
// }

// func (g *GameRoom) StartGame() {
// 	g.GameStarted = true
// }

// func NewGameRoom() *GameRoom {
// 	return &GameRoom{
// 		Game:        nil,
// 		Players:     make([]*Player, 0, MAX_PLAYERS),
// 		ID:          generateUniqueID(),
// 		NumPlayers:  0,
// 		GameStarted: false,
// 		RoomChannel: make(chan string),
// 	}
// }

// func (g *GameRoom) HasSpace() bool {
// 	return g.NumPlayers < 2
// }

// type Player struct {
// 	ID          string
// 	Username    string
// 	Connected   bool
// 	Room        *GameRoom
// 	ChatChannel chan string
// 	Conn        *websocket.Conn
// }

// func generateUniqueID() string {
// 	return "1234567890"
// }

// func NewPlayer(username string, conn *websocket.Conn) *Player {
// 	return &Player{
// 		ID:          generateUniqueID(),
// 		Username:    username,
// 		Connected:   true,
// 		Room:        nil,
// 		ChatChannel: make(chan string),
// 		Conn:        conn,
// 	}
// }
