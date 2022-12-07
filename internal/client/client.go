package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	Servers map[int]string
}

func (c *Client) Run() {
	fmt.Println("SDR - Labo 3 - Client")
	reader := bufio.NewReader(os.Stdin)
	for {
		displayPrompt()
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}

		command, servAddr, err := processInput(input)
		if err != nil {
			log.Println(err)
			continue
		}

		sendCommand(command, servAddr)
	}
}

func displayPrompt() {
	fmt.Println("Available commands:")
	fmt.Println("  - <server number> add <number>")
	fmt.Println("  - <server number> ask")
	fmt.Println("  - <server number> new")
	fmt.Println("  - <server number> stop")
	// TODO: Add a quit command to exit the program and handle ctrl+c
	fmt.Println("Enter a command to send to a connected server on the network:")
}

func processInput(input string) (string, string, error) {
	args := strings.Fields(input)

	if len(args) == 0 {
		return "", "", fmt.Errorf("empty input")
	}
	// TODO: Get server number, then the command with its arguments or return an error

	// TODO: Process server number to get the address in map or return an error

	// TODO: Prepare a json string with corresponding command type or return an error

	return "", "", fmt.Errorf("invalid input")
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

	buffer := make([]byte, 1024)
	n, servAddr, err := connection.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s sent: %s\n", servAddr.String(), string(buffer[0:n]))
}
