package server

type JoinServerMessage struct {
	Username string `json:"username"`
}

type StartGameMessage struct {
	Type    string `json:"type"`
	Context string `json:"context"`
}

type PlayGameMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}
