package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/Lazzzer/labo3-sdr/internal/server"
	"github.com/Lazzzer/labo3-sdr/internal/shared"
)

//go:embed config.json
var config string

func main() {
	if len(os.Args) > 1 {
		log.Fatal("You should not pass any arguments")
	}

	configuration := shared.ParseConfig(config)

	serv := server.Server{Address: configuration.Servers[1], Addresses: configuration.Servers}
	serv.Run()
}
