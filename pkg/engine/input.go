package engine

type InputProvider interface {
	GetPlay(*Table, *Hand, string, *TurnState) Play
}

type TerminalInputProvider struct{}

func (t *TerminalInputProvider) GetPlay(table *Table, hand *Hand, playerName string, turnState *TurnState) Play {

	renderer := NewRenderer(table, hand, playerName, turnState)
	return renderer.UserInputDisplay()
}
