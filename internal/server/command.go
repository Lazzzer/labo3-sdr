// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

// Package serveur propose un serveur UDP connecté dans un réseau de serveurs. Le serveur peut recevoir des commandes de clients UDP et
// créer des élections pour choisir un processus élu avec la charge la moins élevée en suivant l'algorithme de Chang et Roberts.
// Le processus d'élection est géré par des messages inter-processus envoyés par UDP avec acks.
// Les élections resistent à des pannes de processus mais peuvent être perturbées par des cas limites.
package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

// handleCommand gère les commandes reçues des clients UDP.
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
		s.newElectionChan <- true
	case types.Stop:
		os.Exit(1)
	}
	return "Command " + string(command.Type) + " handled", nil
}

// handleAdd gère les ajouts de valeur à la charge du processus.
func (s *Server) handleAdd(command *types.Command) {
	isRunning := <-s.electionStateChan
	if !isRunning {
		s.process.Value += *command.Value
		shared.Log(types.INFO, "New value added to process, value is now: "+strconv.Itoa(s.process.Value))
	} else {
		shared.Log(types.INFO, "Election is running, waiting for election to end")
		<-s.electedChan
		s.process.Value += *command.Value
		shared.Log(types.INFO, "New value added to process, value is now: "+strconv.Itoa(s.process.Value))
	}
}

// handleAsk gère la demande du numéro du processus élu par un client UDP.
func (s *Server) handleAsk() string {
	value := s.getElected()

	var response string
	if value == -1 {
		response = "No election was run"
	} else {
		response = "Process P" + strconv.Itoa(value) + " from Server @" + s.Servers[s.getNextServer(value)] + " was elected"
	}
	shared.Log(types.INFO, "RES TO ASK => "+response)
	return response
}
