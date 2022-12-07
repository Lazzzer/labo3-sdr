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
	Address   string
	Addresses map[int]string
}

func (c *Client) Run() {
	udpAddr, err := net.ResolveUDPAddr("udp4", c.Address)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The UDP server is %s\n", connection.RemoteAddr().String())
	defer connection.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")

		_, err = connection.Write(data)
		if err != nil {
			log.Fatal(err)
		}

		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Exiting UDP client!")
			return
		}

		buffer := make([]byte, 1024)
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Reply: %s\n", string(buffer[0:n]))
	}
}
