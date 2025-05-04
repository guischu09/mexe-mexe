package engine

type InputProvider interface {
	GetPlay(*Table, *Hand, string) Play
}

type TerminalInputProvider struct{}

func (t *TerminalInputProvider) GetPlay(table *Table, hand *Hand, playerName string) Play {

	renderer := NewRenderer(table, hand, playerName)
	return renderer.UserInputDisplay()
}
