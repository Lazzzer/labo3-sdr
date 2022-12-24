// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

// Package serveur propose un serveur UDP connecté dans un réseau de serveurs. Le serveur peut recevoir des commandes de clients UDP et
// créer des élections pour choisir un processus élu avec la charge la moins élevée en suivant l'algorithme de Chang et Roberts.
// Le processus d'élection est géré par des messages inter-processus envoyés par UDP avec acks.
// Les élections resistent à des pannes de processus mais peuvent être perturbées par des cas limites.
package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

// handleMessage gère les messages reçus des autres serveurs lors de l'exécution de l'algorithme de Chang et Roberts.
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

	if s.Debug {
		shared.Log(types.DEBUG, "Throttling message handling")
		time.Sleep(time.Duration(time.Duration(s.DebugDelay) * time.Second))
	}

	switch message.Type {
	case types.Ann:
		s.annChan <- *message
	case types.Res:
		s.resChan <- *message
	default:
		return fmt.Errorf("invalid message type")
	}

	return nil
}

// handleAnn gère les messages de type announcement de l'algorithme de Chang et Roberts.
func (s *Server) handleAnn(message *types.Message) {
	shared.Log(types.MESSAGE, "GOT => Type: announcement, List: "+shared.ShowProcessList(message.Processes, true))
	var messageToSend types.Message
	isProcessInlist := false

	for _, p := range message.Processes {
		if s.process == p {
			isProcessInlist = true
			break
		}
	}

	if isProcessInlist {
		s.elected = shared.GetNbProcessWithMinValue(&message.Processes)
		shared.Log(types.INFO, shared.PURPLE+"Elected process: "+strconv.Itoa(s.elected)+shared.RESET)

		processes := make([]types.Process, 0)
		processes = append(processes, s.process)
		messageToSend = types.Message{Type: types.Res, Elected: s.elected, Processes: processes}
		s.electionState = types.Res
	} else {
		processes := append(message.Processes, s.process)
		messageToSend = types.Message{Type: types.Ann, Elected: -1, Processes: processes}
		s.electionState = types.Ann
	}

	err := s.sendMessage(&messageToSend, s.getNextServer(s.Number))
	if err != nil {
		shared.Log(types.ERROR, "Error while sending message : "+err.Error())
	}
}

// handleRes gère les messages de type result de l'algorithme de Chang et Roberts.
func (s *Server) handleRes(message *types.Message) {
	shared.Log(types.MESSAGE, "GOT => Type: result, Elected: P"+strconv.Itoa(message.Elected)+", List: "+shared.ShowProcessList(message.Processes, false))
	var messageToSend types.Message
	isProcessInlist := false

	for _, p := range message.Processes {
		if s.process.Number == p.Number {
			isProcessInlist = true
			break
		}
	}

	if isProcessInlist {
		return
	}

	if s.electionState == types.Res && s.elected != s.processNumber {
		processes := append(make([]types.Process, 0), s.process)
		messageToSend = types.Message{Type: types.Ann, Elected: -1, Processes: processes}
		s.electionState = types.Ann
	} else if s.electionState == types.Ann {
		s.elected = message.Elected
		shared.Log(types.INFO, shared.PURPLE+"Elected process: "+strconv.Itoa(s.elected)+shared.RESET)
		processes := append(message.Processes, s.process)
		messageToSend = types.Message{Type: types.Res, Elected: s.elected, Processes: processes}
		s.electionState = types.Res
	}

	err := s.sendMessage(&messageToSend, s.getNextServer(s.Number))
	if err != nil {
		shared.Log(types.ERROR, "Error while sending message "+err.Error())
	}
}

// sendMessage envoie un message de type announcement ou result à un autre serveur.
func (s *Server) sendMessage(message *types.Message, destServer int) error {
	if s.Debug {
		shared.Log(types.DEBUG, "Throttling message sending")
		time.Sleep(time.Duration(time.Duration(s.DebugDelay) * time.Second))
	}

	if message.Type == types.Res {
		go func() {
			s.endElectionChan <- true
		}()
	}

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
		stringToLog += shared.ShowProcessList(message.Processes, true)
	} else {
		stringToLog += shared.ShowProcessList(message.Processes, false)
	}

	shared.Log(types.MESSAGE, stringToLog)

	// Wait for acknowledgement from the next process & timeout after 1 second
	buffer := make([]byte, 1024)
	errDeadline := connection.SetReadDeadline(time.Now().Add(time.Duration(s.TimeoutDelay) * time.Second))
	if errDeadline != nil {
		shared.Log(types.ERROR, errDeadline.Error())
	}

	n, _, err := connection.ReadFromUDP(buffer)

	if err != nil {
		if e, ok := err.(net.Error); !ok || e.Timeout() {
			return fmt.Errorf("error while reading from udp: %v", err)
		}
		shared.Log(types.ERROR, "TIMEOUT for ACK from P"+strconv.Itoa(destServer-1))

		err := s.sendMessage(message, s.getNextServer(destServer))
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

// startElection lance une élection en vérifiant qu'aucune élection n'est déjà en cours.
func (s *Server) startElection() {

	if s.electionState == types.Ann {
		shared.Log(types.INFO, "An election is already running")
		return
	}

	shared.Log(types.INFO, shared.PURPLE+"Starting election"+shared.RESET)

	processes := append(make([]types.Process, 0), s.process)
	message := types.Message{Type: types.Ann, Processes: processes}

	err := s.sendMessage(&message, s.getNextServer(s.Number))
	if err != nil {
		shared.Log(types.ERROR, "Error while sending message : "+err.Error())
	}

	s.electionState = types.Ann
}

// getElected retourne le numéro du processus élu.
func (s *Server) getElected() int {
	isRunning := <-s.electionStateChan
	if !isRunning {
		return s.elected
	} else {
		return <-s.electedChan
	}
}

// getNextServer retourne le numéro du prochain serveur à qui envoyer un message.
func (s *Server) getNextServer(current int) int {
	if current == s.nbProcesses {
		return 1
	}
	return current + 1
}
