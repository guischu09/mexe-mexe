package main

import (
	"fmt"
	"mexemexe/internal/client"
	"mexemexe/internal/engine"
	"os"
	"os/signal"
	"syscall"
)

func SigtermListener(signalChan chan os.Signal, stopChan chan struct{}) {
	sig := <-signalChan
	fmt.Printf("Received termination signal (%v), shutting down mexe-mexe game.", sig)

	select {
	case <-stopChan:
	default:
		close(stopChan)
	}
}

func main() {

	serverIP := "192.168.15.6"

	// Instantiate a client
	client := client.NewClient(serverIP, "8888")
	defer client.Close()

	// Set username
	client.SetUsername()

	// Establish websocket connection
	client.SetWebsocketConnection()

	// Send join message - Here should enter authentication
	client.SendJoinMessage()

	// Read join response from server
	client.ReceiveWelcomeMessage()

	// Send start game message to server -- TODO with game options
	client.SendStartGameMessage()

	// Read Join game response from server
	client.ReceiveJoinedGameRoomMessage()

	// Read Game started message from server
	client.ReceiveGameStartedMessage()

	// Set renderer
	renderer := engine.NewRenderer(client.Username)
	client.SetRenderer(renderer)

	// Setup signal and stop channels
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	stopChan := make(chan struct{})

	// Start signal listener
	go SigtermListener(signalChan, stopChan)

	// Start the game
	stopDisplay := make(chan bool)

	// Start the game!
	client.StartGame(stopDisplay)

}
