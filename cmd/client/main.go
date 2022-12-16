package main

import (
	_ "embed"
	"flag"
	"log"
	"os"

	"github.com/Lazzzer/labo3-sdr/internal/client"
	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

//go:embed config.json
var config string

func main() {
	if len(os.Args) > 1 {
		log.Fatal("You should not pass any arguments") // TODO: Maybe useless?
	}

	debug := flag.Bool("debug", false, "Boolean: Run client in debug mode. Default is false")
	flag.Parse()

	configuration, err := shared.Parse[types.Config](config)
	if err != nil {
		log.Fatal(err)
	}

	cl := client.Client{Servers: configuration.Servers}
	cl.Run(*debug)
}
