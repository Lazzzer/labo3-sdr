package main

import (
	_ "embed"
	"flag"
	"log"

	"github.com/Lazzzer/labo3-sdr/internal/client"
	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

//go:embed config.json
var config string
var DEBUG_DELAY = 5  // in seconds
var timeoutValue = 1 // in seconds

func main() {
	if len(flag.Args()) > 1 {
		log.Fatal("usage: go run ./main.go [-debug]")
	}
	//TODO: crash if someone tries to run like this ? "go run ./main.go bite"

	debug := flag.Bool("debug", false, "Boolean: Run client in debug mode. Default is false")
	flag.Parse()

	configuration, err := shared.Parse[types.Config](config)
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		timeoutValue *= DEBUG_DELAY
	}

	cl := client.Client{
		Debug:        *debug,
		Servers:      configuration.Servers,
		TimeoutValue: timeoutValue,
	}
	cl.Run()
}
