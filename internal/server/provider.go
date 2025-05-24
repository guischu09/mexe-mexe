package server

import (
	"log"
	"mexemexe/internal/engine"

	"github.com/gorilla/websocket"
)

type WebsocketOutputProvider struct {
	conn *websocket.Conn
}

func NewWebsocketOutputProvider(conn *websocket.Conn) WebsocketOutputProvider {
	return WebsocketOutputProvider{
		conn: conn,
	}
}

func (w WebsocketOutputProvider) Write(messageType string, data interface{}) {
	log.Printf("DEBUG: Write - Writing message type %s", messageType)
}

type WebsocketInputProvider struct {
	conn *websocket.Conn
}

func NewWebsocketInputProvider(conn *websocket.Conn) WebsocketInputProvider {
	return WebsocketInputProvider{
		conn: conn,
	}
}

func (w WebsocketInputProvider) GetPlay(table engine.Table, hand engine.Hand, playerName string, turnState engine.TurnState) engine.Play {
	log.Printf("DEBUG: GetPlay - Table: %v\r\n", table)
	log.Printf("DEBUG: GetPlay - Hand: %v\r\n", hand)
	log.Printf("DEBUG: GetPlay - TurnState: %v\r\n", turnState)

	gameStateMsg := GameStateMessage{
		Table: table,
		Hand:  hand,
		Turn:  turnState,
	}

	// Send the game state to the client
	err := w.conn.WriteJSON(&gameStateMsg)
	if err != nil {
		log.Printf("error writing to websocket: %v", err)
		log.Println("Closing connection")
		w.conn.Close()
		// Return a "quit" play instead of nil
		return engine.NewQuitPlay("connection_lost")
	}

	// Read the client's move
	var gamePlayMsg GamePlayMessage
	err = w.conn.ReadJSON(&gamePlayMsg)
	if err != nil {
		log.Printf("error reading from websocket: %v", err)
		log.Println("Closing connection")
		w.conn.Close()
		// Return a "quit" play instead of nil
		return engine.NewQuitPlay("connection_lost")
	}

	// Handle case where Play might be nil
	if gamePlayMsg.Play == nil {
		log.Printf("warning: received nil play from client, treating as quit")
		return engine.NewQuitPlay("nil_play")
	}

	return gamePlayMsg.Play
}
