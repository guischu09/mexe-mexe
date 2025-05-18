package server

import "mexemexe/internal/engine"

type JoinServerMessage struct {
	Username string `json:"username"`
}

type MaxCapacityMessage struct {
	Message string `json:"message"`
}

type WelcomeMessage struct {
	Message    string `json:"message"`
	PlayerUUID string `json:"player_uuid"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type StartGameMessage struct {
	Action string `json:"action"`
}

type WaitingRoomMessage struct {
	Message string `json:"message"`
}

type JoinedGameRoomMessage struct {
	Message string `json:"message"`
}

type GameStartedMessage struct {
	Message string `json:"message"`
}

type GameStateMessage struct {
	Table *engine.Table
	Hand  *engine.Hand
	Deck  engine.Deck // Is this needed?
	Turn  *engine.TurnState
}

type GamePlayMessage struct {
	Play engine.Play
}

type GameMessage struct {
	Message string `json:"message"`
}
