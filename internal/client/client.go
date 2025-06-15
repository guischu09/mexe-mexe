package client

import (
	"encoding/json"
	"fmt"
	"log"
	"mexemexe/internal/engine"
	"mexemexe/internal/server"
	"net/url"

	"github.com/gorilla/websocket"
)

// Client defines a connected client
type Client struct {
	ServerIP   string
	ServerPort string
	Renderer   *engine.Renderer
	Username   string
	UUID       string
	Conn       *websocket.Conn
}

// NewClient is Client constructor
func NewClient(ip string, port string) *Client {
	client := Client{
		ServerIP:   ip,
		ServerPort: port,
	}
	return &client
}

// SetUsername sets the username for the client from the user input
func (c *Client) SetUsername() {
	var username string
	fmt.Println("Lets play mexe-mexe!")
	fmt.Println("To join a game room, please enter your username:")
	fmt.Scanf("%s", &username)
	c.Username = username
}

// SetWebsocketConnection establishes a websocket connection to the server
func (c *Client) SetWebsocketConnection() {
	url := url.URL{Scheme: "ws", Host: c.ServerIP + ":" + c.ServerPort, Path: "/ws"}
	log.Printf("connecting to %s\n", url.String())

	// Stablish websocket connection
	ws, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c.Conn = ws
}

func (c *Client) SendJoinMessage() {
	joinMessage := server.JoinServerMessage{
		Username: c.Username,
	}
	err := c.Conn.WriteJSON(joinMessage)
	if err != nil {
		log.Printf("error writing to websocket: %v", err)
		return
	}
}

func (c *Client) ReceiveWelcomeMessage() {
	var welcomeMsg server.WelcomeMessage

	err := c.Conn.ReadJSON(&welcomeMsg)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	c.UUID = welcomeMsg.PlayerUUID
	fmt.Println(welcomeMsg.Message)
}

func (c *Client) SendStartGameMessage() {
	startGameMessage := server.StartGameMessage{
		Action: "start",
	}
	err := c.Conn.WriteJSON(startGameMessage)
	if err != nil {
		log.Printf("error writing to websocket: %v", err)
		return
	}
}

func (c *Client) ReceiveJoinedGameRoomMessage() {
	var joinMsg server.JoinedGameRoomMessage
	err := c.Conn.ReadJSON(&joinMsg)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	fmt.Println(joinMsg.Message)
}

func (c *Client) ReceiveGameStartedMessage() {
	var gameStartedMsg server.GameStartedMessage
	err := c.Conn.ReadJSON(&gameStartedMsg)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	fmt.Println(gameStartedMsg.Message)
}

func (c *Client) SetRenderer(renderer *engine.Renderer) {
	c.Renderer = renderer
}

func (c *Client) ReceiveGameState() server.GameStateMessage {
	var gameState server.GameStateMessage
	err := c.Conn.ReadJSON(&gameState)
	if err != nil {
		log.Fatalf("error reading game state: %v", err)
	}

	fmt.Println("DEBUG: Received game state: \n\r")
	log.Print("DEBUG: Game state: gameState.Hand: \n\r")
	gameState.Hand.Print()
	log.Print("DEBUG: Game state: gameState.Table: \n\r")
	gameState.Table.Print()
	log.Print("DEBUG: Turn state: turnState.HasDrawedCard: \n\r", gameState.Turn.HasDrawedCard)
	log.Print("DEBUG: Turn state: turnState.HasPlayedMeld: \n\r", gameState.Turn.HasPlayedMeld)
	log.Print("DEBUG: Turn state: turnState.PlayerUUID: \n\r", gameState.Turn.PlayerUUID)
	return gameState
}
func (c *Client) ReadFromWebSocket(gameStateChan chan server.GameStateMessage, stopChan chan bool) {
	for {
		gameState := c.ReceiveGameState()
		log.Print("DEBUG: Received game state. \n\r")
		gameStateChan <- gameState
		log.Print("DEBUG: sent game state to channel. \n\r")
		stopChan <- true
	}
}

func (c *Client) StartGame(stopSignal chan bool) {
	gameStateChan := make(chan server.GameStateMessage, 1)
	stopChan := make(chan bool, 1)
	go c.ReadFromWebSocket(gameStateChan, stopChan)

	for {
		fmt.Println("DEBUG: beginning of loop. \n\r")

		// Wait to receive game state from the server
		gameState := <-gameStateChan

		select {
		case <-stopChan:
		default:
		}

		// Determine if it's the player's turn
		freeze := gameState.Turn.PlayerUUID != c.UUID

		c.Renderer.UpdateRenderer(gameState.Table, gameState.Hand, gameState.Turn)

		var play engine.Play
		if freeze {
			play = c.Renderer.DisplayScreen(stopChan)
		} else {
			play = c.Renderer.UserInputDisplay(stopChan)
		}

		// If DisplayScreen returned nil due to stopDisplay signal, just continue the loop
		if play == nil {
			continue
		}

		fmt.Println("DEBUG: Play: \n\r", play)
		fmt.Println("DEBUG: Play type: \n\r", play.GetName())
		fmt.Println("DEBUG: Play cards: \n\r", play.GetCards())

		var gamePlayMsg server.GamePlayMessage
		gamePlayMsg.Play = play

		err := c.Conn.WriteJSON(&gamePlayMsg)
		if err != nil {
			log.Printf("error writing to websocket: %v", err)
			return
		}

		if play.GetName() == engine.QUIT {
			break
		}
	}
}

// Close closes the websocket connection
func (c *Client) Close() {
	c.Conn.Close()
}

type RawGamePlayMessage struct {
	Play json.RawMessage `json:"play"`
}
