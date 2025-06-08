package engine

import (
	"fmt"
	"log"
	"time"

	"mexemexe/internal/service"

	"github.com/gorilla/websocket"
)

type GameStateMessageOut struct {
	Table Table     `json:"table"`
	Hand  Hand      `json:"hand"`
	Turn  TurnState `json:"turn"`
}

type MessageType string

type OutputProvider interface {
	Write(messageType string, data interface{})
	SendState(table Table, hand Hand, turnState TurnState)
}

type TerminalOutputProvider struct{}

func (t TerminalOutputProvider) Write(messageType string, data interface{}) {
	switch messageType {

	case "message":
		fmt.Println(data)
	}
}

type WebsocketOutputProvider struct {
	conn   *websocket.Conn
	logger *service.GameLogger
}

func NewWebsocketOutputProvider(conn *websocket.Conn, logger *service.GameLogger) WebsocketOutputProvider {
	return WebsocketOutputProvider{
		conn:   conn,
		logger: logger,
	}
}

func (w WebsocketOutputProvider) Write(messageType string, data interface{}) {
	log.Printf("DEBUG: Write - Writing message type %s", messageType)
}

func (w WebsocketOutputProvider) SendState(table Table, hand Hand, turnState TurnState) {

	fmt.Println("Sending state to player")
	time.Sleep(5 * time.Second)
	gameState := GameStateMessageOut{
		Table: table,
		Hand:  hand,
		Turn:  turnState,
	}
	err := w.conn.WriteJSON(gameState)
	if err != nil {
		w.logger.Errorf("error writing to websocket: %v", err)
		return
	}
	w.logger.Infof("Successfully sent game state to player")
}
