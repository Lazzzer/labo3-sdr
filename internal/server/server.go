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

	process       types.Process     // Processus courant du serveur
	nbProcesses   int               // Nombre de processus dans le réseau
	processNumber int               // Numéro du processus courant
	electionState types.MessageType // État de l'élection
	elected       int               // Numéro du processus élu

	annChan           chan types.Message // Channel pour les messages d'annonce
	resChan           chan types.Message // Channel pour les messages de réponse
	newElectionChan   chan bool
	electionStateChan chan bool
	endElectionChan   chan bool
	electedChan       chan int
}

func (s *Server) Run() {
	if s.Debug {
		shared.Log(types.DEBUG, "Server started in debug mode")
	}
	s.setup()

	connection := s.startListening()
	defer connection.Close()

	shared.Log(types.INFO, shared.GREEN+"Server #"+strconv.Itoa(s.Number)+" as Process P"+strconv.Itoa(s.process.Number)+" listening on "+s.Address+shared.RESET)

	s.handleCommunications(connection)
}

func (s *Server) setup() {

	s.annChan = make(chan types.Message, 1) // Channel pour les messages d'annonce
	s.resChan = make(chan types.Message, 1) // Channel pour les messages de réponse

	s.newElectionChan = make(chan bool)

	s.electionStateChan = make(chan bool)
	s.endElectionChan = make(chan bool)
	s.electedChan = make(chan int)

	s.nbProcesses = len(s.Servers)
	s.processNumber = s.Number - 1
	s.process = types.Process{Number: s.processNumber, Value: 0}
	s.elected = -1
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
	go func() {
		for {
			select {
			case <-s.newElectionChan:
				s.startElection()
			case <-s.endElectionChan:
			out:
				for {
					select {
					case s.electedChan <- s.elected:
					default:
						break out
					}
				}
			case s.electionStateChan <- s.electionState == types.Ann:
			case message := <-s.annChan:
				s.handleAnn(&message)
			case message := <-s.resChan:
				s.handleRes(&message)
			}
		}
	}()

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
			go func() {
				// Traitement d'une commande si le message n'est pas valide
				response, err := s.handleCommand(communication)
				if err != nil {
					shared.Log(types.ERROR, err.Error())
				}
				// Envoi de la réponse à l'adresse du client
				_, err = connection.WriteToUDP([]byte(response), addr)
				if err != nil {
					shared.Log(types.ERROR, err.Error())
				}
			}()
		}
	}
}
