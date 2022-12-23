// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Lazzzer/labo3-sdr/internal/server"
	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

// TestClient est un client de test
type TestClient struct {
}

var testClient = TestClient{}

var Servers = map[int]string{
	1: "localhost:8091",
	2: "localhost:8092",
	3: "localhost:8093",
}

var srv1 = server.Server{
	Debug:        true,
	DebugDelay:   2,
	Number:       1,
	Address:      Servers[1],
	Servers:      Servers,
	TimeoutDelay: 1,
}

var srv2 = server.Server{
	Debug:        true,
	DebugDelay:   2,
	Number:       2,
	Address:      Servers[2],
	Servers:      Servers,
	TimeoutDelay: 1,
}

var srv3 = server.Server{
	Debug:        true,
	DebugDelay:   2,
	Number:       3,
	Address:      Servers[3],
	Servers:      Servers,
	TimeoutDelay: 1,
}

// init() lance les serveurs de test
func init() {

	println("\nSTEP 0: Start servers\n")
	go srv1.Run()
	go srv3.Run()

	time.Sleep(100 * time.Millisecond)
}

func (tc *TestClient) RunElectionWithDownServer(t *testing.T) {

	println("\nSTEP 1: Add load to servers, 10 to srv1, 5 to srv3 and srv2 is down.\n")

	// add value 10 to srv1
	val1 := 10
	cmd1, err := json.Marshal(types.Command{
		Type:  types.Add,
		Value: &val1,
	})
	if err != nil {
		t.Error(err)
	}

	tc.sendCommand(string(cmd1), &srv1, t)

	// add value 5 to srv3$
	val2 := 5
	cmd2, err := json.Marshal(types.Command{
		Type:  types.Add,
		Value: &val2,
	})
	if err != nil {
		t.Error(err)
	}

	tc.sendCommand(string(cmd2), &srv3, t)

	time.Sleep(100 * time.Millisecond)

	println("\nSTEP 2: Request for election to srv1 then ask for elected to srv1 \n")

	cmdNew, err := json.Marshal(types.Command{
		Type: types.New,
	})
	if err != nil {
		t.Error(err)
	}

	tc.sendCommand(string(cmdNew), &srv1, t)

	time.Sleep(100 * time.Millisecond)

	cmdAsk, err := json.Marshal(types.Command{
		Type: types.Ask,
	})
	if err != nil {
		t.Error(err)
	}

	expected := "Process P2 from Server @localhost:8093 was elected"
	got := tc.sendCommand(string(cmdAsk), &srv1, t)

	if got != expected {
		t.Errorf("\n\nExpected %s, got %s\n\n", expected, got)
	}

	time.Sleep(100 * time.Millisecond)
}

func (tc *TestClient) sendCommand(command string, server *server.Server, t *testing.T) string {
	udpAddr, err := net.ResolveUDPAddr("udp", server.Address)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Error(err)
	}
	defer func(connection *net.UDPConn) {
		err := connection.Close()
		if err != nil {
			shared.Log(types.ERROR, err.Error())
		}
	}(connection)

	_, err = connection.Write([]byte(command + "\n"))
	if err != nil {
		t.Error(err)
	}

	buffer := make([]byte, 1024)
	n, servAddr, err := connection.ReadFromUDP(buffer)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(shared.GREEN + "Server @" + servAddr.String() + " -> " + string(buffer[0:n]) + shared.RESET)

	return string(buffer[0:n])
}

func TestElectionWithDownServers(t *testing.T) {
	testClient.RunElectionWithDownServer(t)
}
