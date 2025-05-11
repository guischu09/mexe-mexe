package engine

import (
	"mexemexe/internal/server"

	"github.com/gorilla/websocket"
)

type InputProvider interface {
	GetPlay(*Table, *Hand, string, *TurnState) Play
}

type TerminalInputProvider struct{}

func (t TerminalInputProvider) GetPlay(table *Table, hand *Hand, playerName string, turnState *TurnState) Play {

	renderer := NewRenderer(table, hand, playerName, turnState)
	return renderer.UserInputDisplay()
}

type WebsocketInputProvider struct {
	conn *websocket.Conn
}

type WebSocketOutputProvider struct {
	conn *websocket.Conn
}

func (w WebsocketInputProvider) GetPlay(table *Table, hand *Hand, playerName string, turnState *TurnState) Play {
	renderer := NewRenderer(table, hand, playerName, turnState)
	return renderer.UserInputDisplay()
}

func (w *WebSocketOutputProvider) Write(msgType string, msg string) {
	payload := server.PlayGameMessage{
		Message: msg,
		Type:    msgType,
	}
	w.conn.WriteJSON(payload)
}
