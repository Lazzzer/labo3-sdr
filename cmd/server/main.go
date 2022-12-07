package main

import (
	_ "embed"
	"flag"
	"log"
	"strconv"

	"github.com/Lazzzer/labo3-sdr/internal/server"
	"github.com/Lazzzer/labo3-sdr/internal/shared"
)

//go:embed config.json
var config string

func main() {
	flag.Parse()
	if flag.Arg(0) == "" {
		log.Fatal("Invalid argument, usage: <server number>")
	}

	number, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatal("Invalid argument, usage: <server number>")
	}

	configuration := shared.ParseConfig(config)

	if number > len(configuration.Servers) || number < 0 {
		log.Fatal("Invalid server number")
	}

	serv := server.Server{Address: configuration.Servers[number], Servers: configuration.Servers}
	serv.Run()
}
