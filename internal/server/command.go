package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func (s *Server) handleCommand(commandStr string) (string, error) {
	command, err := shared.Parse[types.Command](commandStr)
	if err != nil || command.Type == "" {
		return "", fmt.Errorf("invalid command")
	}

	stringToLog := "GOT => Type: " + string(command.Type)

	if command.Type == types.Add {
		stringToLog += ", Value: " + strconv.Itoa(*command.Value)
	}

	shared.Log(types.COMMAND, stringToLog)

	switch command.Type {
	case types.Add:
		s.handleAdd(command)
	case types.Ask:
		return s.handleAsk(), nil
	case types.New:
		s.handleNew()
	case types.Stop:
		os.Exit(1)
	}
	return "Command " + string(command.Type) + " handled", nil
}

func (s *Server) handleAdd(command *types.Command) {
	if electionState == types.Ann {
		// TODO: store for later
	} else {
		process.Value += *command.Value
	}
}

func (s *Server) handleAsk() string {
	if electionState == types.Ann {
		// TODO: wait for election to finish
		return "An election is running"
	}
	return "Process P" + strconv.Itoa(elected) + " from Server @" + s.Servers[getNextServer(elected)] + " was elected"
}

func (s *Server) handleNew() {
	if electionState == types.Ann {
		shared.Log(types.INFO, "An election is already running")
	} else {
		s.startElection()
	}
}
