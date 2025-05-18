package engine

type InputProvider interface {
	GetPlay(*Table, *Hand, string, *TurnState) Play
}

type TerminalInputProvider struct{}

// TODO: Check if this should be a pointer.
func (t TerminalInputProvider) GetPlay(table *Table, hand *Hand, playerName string, turnState *TurnState) Play {

	renderer := NewRenderer(playerName)
	renderer.UpdateRenderer(table, hand, turnState)
	stopSignal := make(chan bool)
	return renderer.UserInputDisplay(stopSignal)
}
