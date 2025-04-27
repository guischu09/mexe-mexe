package engine

import "fmt"

type MessageType string

type OutputProvider interface {
	Write(messageType string, data interface{})
}

type TerminalOutputProvider struct{}

func (t *TerminalOutputProvider) Write(messageType string, data interface{}) {
	switch messageType {

	case "message":
		fmt.Println(data)

	}

}
