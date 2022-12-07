package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Address   string
	Addresses map[int]string
}

func (s *Server) Run() {
	udpAddr, err := net.ResolveUDPAddr("udp4", s.Address)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("UDP server up and listening on address " + s.Address)

	defer connection.Close()
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(addr.String(), " -> ", string(buffer[0:n-1]))

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		data := []byte(strconv.Itoa(101))
		fmt.Printf("data: %s\n", string(data))
		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
