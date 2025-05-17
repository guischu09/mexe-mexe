package main

import (
	"flag"
	"fmt"
	"log"
	"mexemexe/internal/server"
	"net/url"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8888", "http service address")

func main() {

	flag.Parse()
	log.SetFlags(0)

	var username string
	fmt.Println("Lets play mexe-mexe!")
	fmt.Println("To join a game room, please enter your username:")
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

	for {
		var gameMsg server.GameMessage
		err = ws.ReadJSON(&gameMsg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		fmt.Println(gameMsg)
	}

}
