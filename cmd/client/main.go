// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

// Package main est le point d'entrée du programme permettant de démarrer le client.
// Le client a à disposition un fichier de configuration qui contient les adresses des serveurs.
// Il gère aussi un flag debug qui permet d'allonger le temps d'attente des requêtes avant de timeout.
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
var DEBUG_DELAY = 30 // in seconds
var timeoutDelay = 1 // in seconds

// main est la méthode d'entrée du programme
func main() {
	if len(flag.Args()) > 1 {
		log.Fatal("usage: go run ./main.go [-debug]")
	}

	debug := flag.Bool("debug", false, "Boolean: Run client in debug mode. Default is false")
	flag.Parse()

	configuration, err := shared.Parse[types.Config](config)
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		timeoutDelay *= DEBUG_DELAY
	}

	cl := client.Client{
		Debug:        *debug,
		Servers:      configuration.Servers,
		TimeoutDelay: timeoutDelay,
	}
	cl.Run()
}
