package engine

import (
	"encoding/json"
	"mexemexe/internal/service"

	"github.com/gorilla/websocket"
)

type RawGamePlayMessage struct {
	Play json.RawMessage `json:"play"`
}

type InputProvider interface {
	GetPlay(TurnState) Play
	IsConnected() bool
}

// type TerminalInputProvider struct{}

// TODO: Check if this should be a pointer.
// func (t TerminalInputProvider) GetPlay(table *Table, hand Hand, playerName string, turnState TurnState) Play {

// 	renderer := NewRenderer(playerName)
// 	renderer.UpdateRenderer(*table, hand, turnState)
// 	stopSignal := make(chan bool)
// 	return renderer.UserInputDisplay(stopSignal)
// }

type WebsocketInputProvider struct {
	conn   *websocket.Conn
	logger *service.GameLogger
}

func NewWebsocketInputProvider(conn *websocket.Conn, logger *service.GameLogger) *WebsocketInputProvider {
	return &WebsocketInputProvider{
		conn:   conn,
		logger: logger,
	}
}

func (w WebsocketInputProvider) IsConnected() bool {
	return true
}

func (w *WebsocketInputProvider) GetPlay(turnState TurnState) Play {

	var rawMsg RawGamePlayMessage
	err := w.conn.ReadJSON(&rawMsg)
	if err != nil {
		w.logger.Errorf("error reading from websocket: %v", err)
		return NewQuitPlay()
	}

	// Detect play type
	type TypeDetector struct {
		Type string `json:"type"`
	}

	var detector TypeDetector
	err = json.Unmarshal(rawMsg.Play, &detector)
	if err != nil {
		w.logger.Errorf("error detecting play type: %v", err)
		return NewQuitPlay()
	}

	w.logger.Infof("Detected play type: %s", detector.Type)

	// Create concrete Play based on type
	switch detector.Type {
	case "DRAW_CARD":
		return NewDrawCardPlay()
	case "END_TURN":
		return NewEndTurnPlay()
	case "QUIT":
		return NewQuitPlay()
	case "PLAY_MELD":
		var meldPlay MeldPlay
		err = json.Unmarshal(rawMsg.Play, &meldPlay)
		if err != nil {
			w.logger.Errorf("error parsing meld: %v", err)
			return NewQuitPlay()
		}
		return NewMeldPlay(meldPlay.Cards)
	default:
		w.logger.Errorf("unknown play type: %s", detector.Type)
		return NewQuitPlay()
	}
}
