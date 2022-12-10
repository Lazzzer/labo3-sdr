package server

import (
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

var process types.Process // Processus courant du serveur
var nbProcesses int       // Nombre de processus dans le réseau
var processNumber int     // Numéro du processus courant
var value int = 0         // Valeur de la charge du processus courant
var elected *int = nil    // Numéro du processus élu

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

	process = types.Process{Number: processNumber, Value: &value}
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

		response, err := handleMessage(communication)
		if err != nil {
			// Traitement d'une commande si le message n'est pas valide
			response, err = handleCommand(communication)
			if err != nil {
				shared.Log(types.ERROR, err.Error())
				continue
			}
			// Envoi de la réponse à l'adresse du client
			_, err = connection.WriteToUDP([]byte(response), addr)
			if err != nil {
				shared.Log(types.ERROR, err.Error())
			}
			continue
		}
		// Si le message est valide, Envoie de la réponse au processus suivant (au premier si le processus courant est le dernier)
		udpAddr, err := net.ResolveUDPAddr("udp4", s.Servers[s.Number%nbProcesses])
		if err != nil {
			log.Fatal(err)
		}
		_, err = connection.WriteToUDP([]byte(response), udpAddr)
		if err != nil {
			shared.Log(types.ERROR, err.Error())
		}
	}
}
