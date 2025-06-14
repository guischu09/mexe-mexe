package engine

import (
	"fmt"
	"log"

	"mexemexe/internal/service"

	"github.com/gorilla/websocket"
)

var EMPTY_WS_OUTPUT_PROVIDER = WebsocketOutputProvider{
	uuid:   "",
	conn:   nil,
	logger: nil,
}

type GameStateMessageOut struct {
	Table Table     `json:"table"`
	Hand  Hand      `json:"hand"`
	Turn  TurnState `json:"turn"`
}

type MessageType string

type OutputProvider interface {
	Write(messageType string, data interface{})
	SendState(table Table, hand Hand, turnState TurnState)
	GetUUID() string
}

type TerminalOutputProvider struct{}

func (t TerminalOutputProvider) Write(messageType string, data interface{}) {
	switch messageType {

	case "message":
		fmt.Println(data)
	}
}

type WebsocketOutputProvider struct {
	uuid   string
	conn   *websocket.Conn
	logger *service.GameLogger
}

func NewWebsocketOutputProvider(conn *websocket.Conn, uuid string, logger *service.GameLogger) WebsocketOutputProvider {
	return WebsocketOutputProvider{
		uuid:   uuid,
		conn:   conn,
		logger: logger,
	}
}

func (w WebsocketOutputProvider) GetUUID() string {
	return w.uuid
}

func (w WebsocketOutputProvider) Write(messageType string, data interface{}) {
	log.Printf("DEBUG: Write - Writing message type %s", messageType)
}

func (w WebsocketOutputProvider) SendState(table Table, hand Hand, turnState TurnState) {

	log.Printf("DEBUG: SendState - Sending state to player %s", w.uuid)
	// time.Sleep(5 * time.Second)
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
