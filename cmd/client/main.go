package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mexemexe/internal/engine"
	"mexemexe/internal/server"
	"net/url"

	"github.com/gorilla/websocket"
)

type RawGamePlayMessage struct {
	Play json.RawMessage `json:"play"`
}

var addr = flag.String("addr", "192.168.15.6:8888", "http service address")

func main() {

	flag.Parse()
	log.SetFlags(0)

	var username string
	fmt.Println("Lets play mexe-mexe!")
	fmt.Println("To join a game room, please enter your username:")
	// username = "guilherme"
	fmt.Scanf("%s", &username)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s\n", u.String())

	// Stablish websocket connection
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	// Send join message - Here should enter authentication
	joinMessage := server.JoinServerMessage{
		Username: username,
	}
	err = ws.WriteJSON(joinMessage)
	if err != nil {
		log.Printf("error writing to websocket: %v", err)
		return
	}
	// Read join response from server
	var welcomeMsg server.WelcomeMessage

	err = ws.ReadJSON(&welcomeMsg)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	playerUUID := welcomeMsg.PlayerUUID
	fmt.Println(welcomeMsg.Message)

	// Send start game message to server -- TODO with game options
	startGameMessage := server.StartGameMessage{
		Action: "start",
	}
	err = ws.WriteJSON(startGameMessage)
	if err != nil {
		log.Printf("error writing to websocket: %v", err)
		return
	}

	// Read Join game response from server
	var joinMsg server.JoinedGameRoomMessage
	err = ws.ReadJSON(&joinMsg)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	fmt.Println(joinMsg.Message)

	// Read Game started message from server
	var gameStartedMsg server.GameStartedMessage
	err = ws.ReadJSON(&gameStartedMsg)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	fmt.Println(gameStartedMsg.Message)

	renderer := engine.NewRenderer(username)
	stopDisplay := make(chan bool)

	// Game started! Handle the request-response pattern
	for {
		fmt.Println("DEBUG: beggining of loop.")
		// Wait for game state from server
		var gameState server.GameStateMessage
		err := ws.ReadJSON(&gameState)
		if err != nil {
			log.Printf("error reading game state: %v", err)
			return
		}

		fmt.Println("DEBUG: Received game state: ")
		// fmt.Println("DEBUG: Game state: ", gameState)
		log.Print("player :: !> DEBUG: Game state: gameState.Hand: ")
		gameState.Hand.Print()
		log.Print("player :: !> DEBUG: Game state: gameState.Table: ")
		gameState.Table.Print()
		log.Print("player :: !> DEBUG: Turn state: turnState.HasDrawedCard: ", gameState.Turn.HasDrawedCard)
		log.Print("player :: !> DEBUG: Turn state: turnState.HasPlayedMeld: ", gameState.Turn.HasPlayedMeld)
		log.Print("player :: !> DEBUG: Turn state: turnState.PlayerUUID: ", gameState.Turn.PlayerUUID)

		// Stop any existing display
		select {
		case stopDisplay <- true:
		default:
		}

		// Determine if it's the player's turn
		freeze := gameState.Turn.PlayerUUID != playerUUID

		// fmt.Println("This player uuid: ", playerUUID)
		// fmt.Println("Turn player uuid: ", gameState.Turn.PlayerUUID)

		renderer.UpdateRenderer(gameState.Table, gameState.Hand, gameState.Turn)

		// fmt.Println("Freeze: ", freeze)

		var play engine.Play
		if freeze {
			play = renderer.DisplayScreen(stopDisplay)

		} else {
			play = renderer.UserInputDisplay(stopDisplay)
		}

		fmt.Println("Play: ", play)
		fmt.Println("Play type: ", play.GetName())
		fmt.Println("Play cards: ", play.GetCards())

		// Send the play back to server immediately
		if play != nil {
			var gamePlayMsg server.GamePlayMessage
			gamePlayMsg.Play = play

			err := ws.WriteJSON(&gamePlayMsg)
			if err != nil {
				log.Printf("error writing to websocket: %v", err)
				return
			}

			// If it's a quit play, break the loop
			if play.GetName() == engine.QUIT {
				break
			}
		}
	}
}
