package main

import (
	_ "embed"
	"github.com/Lazzzer/labo3-sdr/internal/client"
	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"log"
	"os"
)

//embed config.json
var config string

func main() {
	if len(os.Args) > 1 {
		log.Fatal("You should not pass any arguments")
	}

	configuration := shared.Parse(config)

	cl := client.Client{Address: "localhost", Addresses: configuration.Addresses}
	cl.Run()
}
