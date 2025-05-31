package server

import (
	"encoding/json"
	"log"
	"mexemexe/internal/engine"
	"mexemexe/internal/service"

	"github.com/gorilla/websocket"
)

type WebsocketOutputProvider struct {
	conn   *websocket.Conn
	logger service.GameLogger
}

func NewWebsocketOutputProvider(conn *websocket.Conn, logger service.GameLogger) WebsocketOutputProvider {
	return WebsocketOutputProvider{
		conn:   conn,
		logger: logger,
	}
}

func (w WebsocketOutputProvider) Write(messageType string, data interface{}) {
	log.Printf("DEBUG: Write - Writing message type %s", messageType)
}

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

func (w *WebsocketInputProvider) GetPlay(table engine.Table, hand engine.Hand, playerName string, turnState engine.TurnState) engine.Play {
	type RawGamePlayMessage struct {
		Play json.RawMessage `json:"play"`
	}

	var rawMsg RawGamePlayMessage
	err := w.conn.ReadJSON(&rawMsg)
	if err != nil {
		w.logger.Errorf("error reading from websocket: %v", err)
		return engine.NewQuitPlay()
	}

	// Detect play type
	type TypeDetector struct {
		Type string `json:"type"`
	}

	var detector TypeDetector
	err = json.Unmarshal(rawMsg.Play, &detector)
	if err != nil {
		w.logger.Errorf("error detecting play type: %v", err)
		return engine.NewQuitPlay()
	}

	w.logger.Infof("Detected play type: %s", detector.Type)

	// Create concrete Play based on type
	switch detector.Type {
	case "DRAW_CARD":
		return engine.NewDrawCardPlay()
	case "END_TURN":
		return engine.NewEndTurnPlay()
	case "QUIT":
		return engine.NewQuitPlay()
	case "PLAY_MELD":
		var meldData engine.MeldPlay
		err = json.Unmarshal(rawMsg.Play, &meldData)
		if err != nil {
			w.logger.Errorf("error parsing meld: %v", err)
			return engine.NewQuitPlay()
		}
		return meldData
	default:
		w.logger.Errorf("unknown play type: %s", detector.Type)
		return engine.NewQuitPlay()
	}
}
