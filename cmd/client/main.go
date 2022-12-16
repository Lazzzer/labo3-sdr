package main

import (
	_ "embed"
	"flag"
	"github.com/Lazzzer/labo3-sdr/internal/client"
	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
	"log"
)

//go:embed config.json
var config string

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

	cl := client.Client{Servers: configuration.Servers}
	cl.Run(*debug)
}
