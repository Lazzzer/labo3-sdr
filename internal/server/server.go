package server

import (
	"log"
	"net"
	"strconv"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

type Server struct {
	Debug        bool           // Mode debug
	DebugDelay   int            // Valeur du délai de debug
	Number       int            // Numéro du serveur
	Address      string         // Adresse du serveur
	Servers      map[int]string // Map des serveurs
	TimeoutDelay int            // Valeur du timeout
}

var process types.Process           // Processus courant du serveur
var nbProcesses int                 // Nombre de processus dans le réseau
var processNumber int               // Numéro du processus courant
var electionState types.MessageType // État de l'élection
var elected int = -1                // Numéro du processus élu

func (s *Server) Run() {
	if s.Debug {
		shared.Log(types.DEBUG, "Server started in debug mode")
	}
	s.setupProcessValue()

	connection := s.startListening()
	defer connection.Close()

	shared.Log(types.INFO, shared.GREEN+"Server #"+strconv.Itoa(s.Number)+" as Process P"+strconv.Itoa(process.Number)+" listening on "+s.Address+shared.RESET)

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

		communication := string(buffer[0:n])
		err = s.handleMessage(connection, addr, communication)
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
