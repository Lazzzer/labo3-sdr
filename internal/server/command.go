package server

import (
	"fmt"
	"os"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func handleCommand(commandStr string) (string, error) {
	command, err := shared.Parse[types.Command](commandStr)
	if err != nil || command.Type == "" {
		return "", fmt.Errorf("invalid command")
	}

	switch command.Type {
	case types.Add:
		handleAdd(command)
	case types.Ask:
		handleAsk()
	case types.New:
		handleNew()
	case types.Stop:
		os.Exit(1)
	}
	return "Command " + string(command.Type) + " handled", nil
}

func handleAdd(command *types.Command) {
	// TODO: handle add command
}

func handleAsk() {
	// TODO: handle ask command
}

func handleNew() {
	// TODO: handle new command
}
