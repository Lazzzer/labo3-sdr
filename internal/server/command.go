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
		go s.handleAdd(command)
	case types.Ask:
		return s.handleAsk(), nil
	case types.New:
		newElectionChan <- true
	case types.Stop:
		os.Exit(1)
	}
	return "Command " + string(command.Type) + " handled", nil
}

func (s *Server) handleAdd(command *types.Command) {
	isRunning := <-electionStateChan
	if !isRunning {
		process.Value += *command.Value
		shared.Log(types.INFO, "New value added to process, value is now: "+strconv.Itoa(process.Value))
	} else {
		shared.Log(types.INFO, "Election is running, waiting for election to end")
		<-electedChan
		process.Value += *command.Value
		shared.Log(types.INFO, "New value added to process, value is now: "+strconv.Itoa(process.Value))
	}
}

func (s *Server) handleAsk() string {
	value := getElected()
	response := "Process P" + strconv.Itoa(value) + " from Server @" + s.Servers[getNextServer(value)] + " was elected"
	shared.Log(types.INFO, "RES TO ASK => "+response)
	return response
}
