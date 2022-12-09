package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

type Client struct {
	Servers map[int]string
}

var running = true

func (c *Client) Run() {
	fmt.Println("SDR - Labo 3 - Client")
	reader := bufio.NewReader(os.Stdin)
	for running {
		displayPrompt()
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}

		command, servAddr, err := processInput(input, c)
		if err != nil {
			log.Println(err)
			continue
		}

		sendCommand(command, servAddr)
	}
	fmt.Println("Good bye!")
}

func displayPrompt() {
	fmt.Println("Available commands:")
	fmt.Println("  - <server number> add <number>")
	fmt.Println("  - <server number> ask")
	fmt.Println("  - <server number> new")
	fmt.Println("  - <server number> stop")
	// TODO: handle ctrl+c
	fmt.Println("  - quit")
	fmt.Println("Enter a command to send to a connected server on the network:")
}

func processInput(input string, c *Client) (string, string, error) {
	args := strings.Fields(input)

	// String vide
	if len(args) == 0 {
		return "", "", fmt.Errorf("empty input")
	}

	// Quit
	if args[0] == string(types.Quit) {
		running = false
		return "", "", fmt.Errorf("client quitting program")
		// TODO: Refactor, should not return an error, we can reuse ctrl+c logic with a signal channel
	}

	if len(args) < 2 {
		return "", "", fmt.Errorf("missing command")
	}

	// Vérification du numéro du serveur
	srvNumber, err := strconv.Atoi(args[0])
	if err != nil || srvNumber < 1 || srvNumber > len(c.Servers) {
		return "", "", fmt.Errorf("invalid server number")
	}

	// Vérification de la commande
	command := types.Command{Value: nil}
	switch args[1] {
	case string(types.Add):
		if len(args) != 3 {
			return "", "", fmt.Errorf("invalid add command")
		}
		value, err := strconv.Atoi(args[2])
		if err != nil || value <= 0 {
			return "", "", fmt.Errorf("invalid add command")
		}
		command.Type = types.Add
		command.Value = &value
	case string(types.Ask):
		command.Type = types.Ask
	case string(types.New):
		command.Type = types.New
	case string(types.Stop):
		command.Type = types.Stop
	default:
		return "", "", fmt.Errorf("unknown command")
	}

	// Création du json
	if jsonCommand, err := json.Marshal(command); err == nil {
		return string(jsonCommand), c.Servers[srvNumber], nil
	} else {
		return "", "", fmt.Errorf("invalid command")
	}
}

func sendCommand(command string, address string) {

	udpAddr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	_, err = connection.Write([]byte(command + "\n"))
	if err != nil {
		log.Fatal(err)
	}

	// TODO : handle timeout (>= 1s means server is down)
	buffer := make([]byte, 1024)
	n, servAddr, err := connection.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s sent: %s\n", servAddr.String(), string(buffer[0:n]))
}
