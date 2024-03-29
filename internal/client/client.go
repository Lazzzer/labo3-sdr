// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

// Package client propose un client UDP envoyant des commandes sous forme de string json à des serveurs du réseau.
//
// Le client parse l'entrée de l'utilisateur et envoie la commande correspondante au serveur. Si le serveur destinataire
// ne répond pas dans un délai défini par  TimeoutDelay seconde, le client affiche un message d'erreur.
package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

// Client est la structure représentant le client UDP capable de faire des commandes aux serveurs connectés.
type Client struct {
	Debug        bool           // Mode debug
	Servers      map[int]string // Map des serveurs accessibles
	TimeoutDelay int            // Valeur du timeout
}

var exitChan = make(chan os.Signal, 1) // Channel qui gère à la capture du CTRL+C

// Messages d'erreur
var invalidCommand = "invalid command"
var wrongServerNumber = "invalid server number"
var emptyInput = "empty input"
var chargeMustBePositive = "charge must be a positive integer"

// Run est la méthode principale du client. Elle gère l'entrée de l'utilisateur et envoie les commandes aux serveurs.
func (c *Client) Run() {
	signal.Notify(exitChan, syscall.SIGINT)

	if c.Debug {
		shared.Log(types.DEBUG, "Client started in debug mode")
	}

	go func() {
		<-exitChan
		fmt.Println("Good bye!")
		os.Exit(0)
	}()

	fmt.Println("SDR - Labo 3 - Client")
	reader := bufio.NewReader(os.Stdin)
	for {
		displayPrompt()
		input, err := reader.ReadString('\n')
		if err != nil {
			shared.Log(types.ERROR, err.Error())
			continue
		}

		command, servAddr, err := processInput(input, c)
		if err != nil {
			shared.Log(types.ERROR, err.Error())
			continue
		}

		c.sendCommand(command, servAddr)
	}
}

// displayPrompt affiche les commandes disponibles pour l'utilisateur.
func displayPrompt() {
	fmt.Println("\nAvailable commands:")
	fmt.Println("  - <server number> add <number>")
	fmt.Println("  - <server number> ask")
	fmt.Println("  - <server number> new")
	fmt.Println("  - <server number> stop")
	fmt.Println("  - quit")
	fmt.Println("Enter a command to send to a connected server on the network:")
}

// processInput parse l'entrée de l'utilisateur et retourne un tuple contenant :
// - la commande à envoyer au serveur sous forme de string json
// - l'adresse du serveur auquel envoyer la commande
// - une erreur si l'entrée est invalide
func processInput(input string, c *Client) (string, string, error) {
	args := strings.Fields(input)

	// String vide
	if len(args) == 0 {
		return "", "", fmt.Errorf(emptyInput)
	}

	// Quit
	if args[0] == string(types.Quit) {
		exitChan <- syscall.SIGINT
	}

	if len(args) < 2 || len(args) > 3 {
		return "", "", fmt.Errorf(invalidCommand)
	}

	// Vérification du numéro du serveur
	srvNumber, err := strconv.Atoi(args[0])
	if err != nil || srvNumber < 1 || srvNumber > len(c.Servers) {
		return "", "", fmt.Errorf(wrongServerNumber)
	}

	// Vérification des commandes
	command := types.Command{Value: nil}
	if len(args) == 3 && args[1] == string(types.Add) { // Vérification de la commande ask
		value, err := strconv.Atoi(args[2])
		if err == nil && value >= 0 {
			command.Type = types.Add
			command.Value = &value
		} else {
			return "", "", fmt.Errorf(chargeMustBePositive)
		}
	} else if len(args) == 2 { // Vérification des autres commandes
		switch args[1] {
		case string(types.Ask):
			command.Type = types.Ask
		case string(types.New):
			command.Type = types.New
		case string(types.Stop):
			command.Type = types.Stop
		default:
			return "", "", fmt.Errorf(invalidCommand)
		}
	} else {
		return "", "", fmt.Errorf(invalidCommand)
	}

	// Création du json
	if jsonCommand, err := json.Marshal(command); err == nil {
		return string(jsonCommand), c.Servers[srvNumber], nil
	} else {
		return "", "", fmt.Errorf("error while marshalling command")
	}
}

// sendCommand envoie une commande au serveur spécifié. Elle s'occupe de la connexion UDP et de la fermeture de celle-ci.
// Elle a également la possibilité de timeout si le serveur ne répond pas.
func (c *Client) sendCommand(command string, address string) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		shared.Log(types.ERROR, err.Error())
		return
	}
	defer func(connection *net.UDPConn) {
		err := connection.Close()
		if err != nil {
			shared.Log(types.ERROR, err.Error())
		}
	}(connection)

	_, err = connection.Write([]byte(command + "\n"))
	if err != nil {
		shared.Log(types.ERROR, err.Error())
		return
	}

	buffer := make([]byte, 1024)
	errDeadLine := connection.SetReadDeadline(time.Now().Add(time.Duration(c.TimeoutDelay) * time.Second))
	if errDeadLine != nil {
		return
	}
	n, servAddr, err := connection.ReadFromUDP(buffer)

	if err != nil {
		if e, ok := err.(net.Error); !ok || !e.Timeout() {
			fmt.Println(shared.RED + "Error while reading from server @" + udpAddr.String() + shared.RESET)
			return
		}
		fmt.Println(shared.RED + "Server @" + udpAddr.String() + " is unreachable" + shared.RESET)
		return
	}

	fmt.Println(shared.GREEN + "Server @" + servAddr.String() + " -> " + string(buffer[0:n]) + shared.RESET)
}
