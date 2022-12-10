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

	shared.Log(types.INFO, shared.GREEN+"Server #"+strconv.Itoa(s.Number)+" listening on "+s.Address+shared.RESET)
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
			response, err = handleCommand(communication)
			if err != nil {
				shared.Log(types.ERROR, err.Error())
				continue
			}
		}

		_, err = connection.WriteToUDP([]byte(response), addr)
		if err != nil {
			shared.Log(types.ERROR, err.Error())
		}
	}
}
