package server

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strconv"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func (s *Server) handleMessage(messageStr string) error {
	message, err := shared.Parse[types.Message](messageStr)
	if err != nil || message.Type == "" {
		return fmt.Errorf("invalid message type")
	}

	switch message.Type {
	case types.Ann:
		shared.Log(types.MESSAGE, "GOT => Type: announcement, List: "+showProcessList(message.Processes, true))
		s.handleAnn(message)
	case types.Res:
		shared.Log(types.MESSAGE, "GOT => Type: result, Elected: P"+strconv.Itoa(message.Elected)+", List: "+showProcessList(message.Processes, false))
		s.handleRes(message)
	default:
		return fmt.Errorf("invalid message type")
	}

	return nil
}

func (s *Server) handleAnn(message *types.Message) {
	var messageToSend types.Message
	isProcessInlist := false

	for _, p := range message.Processes {
		if process == p {
			isProcessInlist = true
			break
		}
	}

	if isProcessInlist {
		elected = getNbProcessWithMinValue(&message.Processes)
		shared.Log(types.INFO, shared.PINK+"Elected process: "+strconv.Itoa(elected)+shared.RESET)

		processes := make([]types.Process, 0)
		processes = append(processes, process)
		messageToSend = types.Message{Type: types.Res, Elected: elected, Processes: processes}
		electionState = types.Res
	} else {
		processes := append(message.Processes, process)
		messageToSend = types.Message{Type: types.Ann, Elected: -1, Processes: processes}
		electionState = types.Ann
	}

	err := s.sendMessage(&messageToSend)
	if err != nil {
		shared.Log(types.ERROR, "Error while sending message")
	}
}

func (s *Server) handleRes(message *types.Message) {
	var messageToSend types.Message
	isProcessInlist := false

	for _, p := range message.Processes {
		if process == p {
			isProcessInlist = true
			break
		}
	}

	if isProcessInlist {
		return
	}

	if electionState == types.Res && elected != processNumber {
		processes := append(make([]types.Process, 0), process)
		messageToSend = types.Message{Type: types.Ann, Elected: -1, Processes: processes}
		electionState = types.Ann
	} else if electionState == types.Ann {
		elected = message.Elected
		shared.Log(types.INFO, shared.PINK+"Elected process: "+strconv.Itoa(elected)+shared.RESET)
		processes := append(message.Processes, process)
		messageToSend = types.Message{Type: types.Res, Elected: elected, Processes: processes}
	}

	err := s.sendMessage(&messageToSend)
	if err != nil {
		shared.Log(types.ERROR, "Error while sending message")
	}
}

func (s *Server) sendMessage(message *types.Message) error {

	messageJson, err := json.Marshal(message)
	if err != nil {
		shared.Log(types.ERROR, err.Error())
		return err
	}

	destServer := s.Number + 1
	if destServer > nbProcesses {
		destServer = 1
	}

	destUdpAddr, err := net.ResolveUDPAddr("udp4", s.Servers[destServer])
	if err != nil {
		return err
	}
	connection, err := net.DialUDP("udp", nil, destUdpAddr)
	if err != nil {
		return err
	}
	_, err = connection.Write(messageJson)
	if err != nil {
		return err
	}
	// TODO : better log message sent

	stringToLog := "SENT TO P" + strconv.Itoa(destServer-1) + " => Type: " + string(message.Type) + ", List: "

	if message.Type == types.Ann {
		stringToLog += showProcessList(message.Processes, true)
	} else {
		stringToLog += showProcessList(message.Processes, false)
	}

	shared.Log(types.MESSAGE, stringToLog)
	return nil
}

func (s *Server) startElection() {
	shared.Log(types.INFO, shared.PINK+"Starting election"+shared.RESET)

	processes := append(make([]types.Process, 0), process)
	message := types.Message{Type: types.Ann, Processes: processes}

	s.sendMessage(&message)
	electionState = types.Ann
}

// TODO: move to utils
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

// TODO: move to utils
func showProcessList(processes []types.Process, withValue bool) string {
	var list string
	list = "["
	for i, p := range processes {
		list += "P" + strconv.Itoa(p.Number)
		if withValue {
			list += ":" + strconv.Itoa(p.Value)
		}
		if i != len(processes)-1 {
			list += ", "
		}
	}
	list += "]"
	return list
}
