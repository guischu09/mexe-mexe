package server

type JoinServerMessage struct {
	Username string `json:"username"`
}

type MaxCapacityMessage struct {
	Message string `json:"message"`
}

type WelcomeMessage struct {
	Message string `json:"message"`
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

type GameMessage struct {
	Message string `json:"message"`
}
