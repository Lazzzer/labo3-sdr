package server

import (
	"fmt"
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

func handleCommand(commandStr string) (string, error) {
	command, err := shared.Parse[types.Command](commandStr)
	if err != nil || command.Type == "" {
		return "", fmt.Errorf("invalid command")
	}

	log.Println(command)

	return "command", nil
}

func handleMessage(messageStr string) (string, error) {
	message, err := shared.Parse[types.Message](messageStr)
	if err != nil || message.Type == "" {
		return "", fmt.Errorf("invalid message")
	}

	log.Println(message)

	return "message", nil
}

func (s *Server) Run() {
	// value := 0

	udpAddr, err := net.ResolveUDPAddr("udp4", s.Address)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	log.Printf("Server #" + strconv.Itoa(s.Number) + " listening on address " + s.Address)

	buffer := make([]byte, 1024)

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}

		communication := string(buffer[0 : n-1])
		log.Println(addr.String(), " -> ", communication)

		response, err := handleMessage(communication)
		if err != nil {
			response, err = handleCommand(communication)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Printf("data: %s\n", string(response))
		_, err = connection.WriteToUDP([]byte(response), addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
