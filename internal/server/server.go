package server

import (
	"encoding/json"
	"log"
	"net"
	"strconv"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

type Server struct {
	Number  int
	Address string
	Servers map[int]string
}

var process types.Process           // Processus courant du serveur
var nbProcesses int                 // Nombre de processus dans le réseau
var processNumber int               // Numéro du processus courant
var electionState types.MessageType // État de l'élection
var elected int = -1                // Numéro du processus élu

func (s *Server) Run() {
	s.setupProcessValue()

	connection := s.startListening()
	defer connection.Close()

	shared.Log(types.INFO, shared.GREEN+"Server #"+strconv.Itoa(s.Number)+" listening on "+s.Address+shared.RESET)

	s.handleCommunications(connection)
}

func (s *Server) setupProcessValue() {
	nbProcesses = len(s.Servers)

	processNumber = s.Number - 1

	process = types.Process{Number: processNumber, Value: 0}
}

func (s *Server) startListening() *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp4", s.Address)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	return connection
}

func (s *Server) handleCommunications(connection *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			shared.Log(types.ERROR, err.Error())
			continue
		}

		communication := string(buffer[0 : n-1])
		shared.Log(types.INFO, shared.YELLOW+addr.String()+" -> "+communication+shared.RESET)

		err = s.handleMessage(communication)
		if err != nil {
			// Traitement d'une commande si le message n'est pas valide
			response, err := s.handleCommand(communication)
			if err != nil {
				shared.Log(types.ERROR, err.Error())
				continue
			}
			// Envoi de la réponse à l'adresse du client
			_, err = connection.WriteToUDP([]byte(response), addr)
			if err != nil {
				shared.Log(types.ERROR, err.Error())
			}
		}
	}
}

func (s *Server) sendMessage(message string) error {

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
	_, err = connection.Write([]byte(message + "\n"))
	if err != nil {
		return err
	}
	shared.Log(types.INFO, shared.CYAN+"Process "+strconv.Itoa(process.Number)+" to Process "+strconv.Itoa(destServer-1)+" => "+string(message)+shared.RESET)
	return nil
}

func (s *Server) startElection() {
	shared.Log(types.INFO, shared.PINK+"Starting election"+shared.RESET)

	processes := make([]types.Process, 0)
	processes = append(processes, process)
	message := types.Message{Type: types.Ann, Processes: processes}

	messageJson, err := json.Marshal(message)
	if err != nil {
		shared.Log(types.ERROR, err.Error())
		return
	}

	s.sendMessage(string(messageJson))
	electionState = types.Ann
}
