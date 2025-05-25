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
	Table engine.Table     `json:"table"`
	Hand  engine.Hand      `json:"hand"`
	Turn  engine.TurnState `json:"turn"`
}

type GamePlayMessage struct {
	Play engine.Play `json:"play"`
}

type GameMessage struct {
	Message string `json:"message"`
}
