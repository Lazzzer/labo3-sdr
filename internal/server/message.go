package server

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func (s *Server) handleMessage(connection *net.UDPConn, addr *net.UDPAddr, messageStr string) error {
	message, err := shared.Parse[types.Message](messageStr)
	if err != nil || message.Type == "" {
		return fmt.Errorf("invalid message type")
	}

	// Send acknoledgement to message sender
	responseJson, err := json.Marshal(types.Acknowledgement{From: s.Number})
	if err != nil {
		return err
	}
	_, err = connection.WriteToUDP([]byte(responseJson), addr)
	if err != nil {
		return err
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

	err := s.sendMessage(&messageToSend, getNextServer(s.Number))
	if err != nil {
		shared.Log(types.ERROR, "Error while sending message : "+err.Error())
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
		electionState = types.Res
	}

	err := s.sendMessage(&messageToSend, getNextServer(s.Number))
	if err != nil {
		shared.Log(types.ERROR, "Error while sending message "+err.Error())
	}
}

func (s *Server) sendMessage(message *types.Message, destServer int) error {

	messageJson, err := json.Marshal(message)
	if err != nil {
		shared.Log(types.ERROR, err.Error())
		return err
	}

	destUdpAddr, err := net.ResolveUDPAddr("udp", s.Servers[destServer])
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

	stringToLog := "SENT TO P" + strconv.Itoa(destServer-1) + " => Type: " + string(message.Type) + ", List: "

	if message.Type == types.Ann {
		stringToLog += showProcessList(message.Processes, true)
	} else {
		stringToLog += showProcessList(message.Processes, false)
	}

	shared.Log(types.MESSAGE, stringToLog)

	// Wait for acknowledgement from the next process & timeout after 1 second
	buffer := make([]byte, 1024)
	errDeadline := connection.SetReadDeadline(time.Now().Add(5 * time.Second))
	if errDeadline != nil {
		shared.Log(types.ERROR, errDeadline.Error())
	}

	n, _, err := connection.ReadFromUDP(buffer)

	if err != nil {
		if e, ok := err.(net.Error); !ok || e.Timeout() {
			return fmt.Errorf("error while reading from udp: %v", err)
		}
		shared.Log(types.ERROR, "TIMEOUT for ACK from P"+strconv.Itoa(destServer-1))

		err := s.sendMessage(message, getNextServer(destServer))
		if err != nil {
			shared.Log(types.ERROR, "Error while sending message : "+err.Error())
		}
		return nil
	}

	messageAck := string(buffer[0:n])
	ack, err := shared.Parse[types.Acknowledgement](messageAck)
	if err != nil {
		return err
	}

	if ack.From != destServer {
		return fmt.Errorf("ack from wrong process")
	}

	return nil
}

func (s *Server) startElection() {
	shared.Log(types.INFO, shared.PINK+"Starting election"+shared.RESET)

	processes := append(make([]types.Process, 0), process)
	message := types.Message{Type: types.Ann, Processes: processes}

	err := s.sendMessage(&message, getNextServer(s.Number))
	if err != nil {
		shared.Log(types.ERROR, "Error while sending message : "+err.Error())
	}

	electionState = types.Ann
}

func getNextServer(current int) int {
	if current == nbProcesses {
		return 1
	}
	return current + 1
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
