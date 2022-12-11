package server

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func (s *Server) handleMessage(messageStr string) error {
	message, err := shared.Parse[types.Message](messageStr)
	if err != nil || message.Type == "" {
		return fmt.Errorf("invalid message")
	}

	switch message.Type {
	case types.Ann:
		s.handleAnn(message)
	case types.Res:
		handleRes(message)
	}

	return nil
}

func (s *Server) handleAnn(message *types.Message) {
	var messageToSend types.Message
	isInList := false

	for _, p := range message.Processes {
		if process == p {
			isInList = true
			break
		}
	}

	if isInList {
		elected = getNbProcessWithMinValue(&message.Processes)

		processes := make([]types.Process, 0)
		processes = append(processes, process)
		messageToSend = types.Message{Type: types.Res, Elected: elected, Processes: processes}
		electionState = types.Res
	} else {
		processes := append(message.Processes, process)
		messageToSend = types.Message{Type: types.Ann, Processes: processes}
		electionState = types.Ann
	}

	messageJson, err := json.Marshal(messageToSend)
	if err != nil {
		shared.Log(types.ERROR, "Error while marshalling message")
	}

	err = s.sendMessage(string(messageJson))
	if err != nil {
		shared.Log(types.ERROR, "Error while sending message")
	}
}

func getNbProcessWithMinValue(processes *[]types.Process) int {
	minValue := math.MaxInt
	minProcessNumber := -1

	for _, p := range *processes {
		if p.Value < minValue {
			minValue = p.Value
			minProcessNumber = p.Number
		}
	}

	return minProcessNumber
}

func handleRes(message *types.Message) {
	// TODO: handle res message
}
