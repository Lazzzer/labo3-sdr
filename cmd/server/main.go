// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

// Package main est le point d'entrée du programme permettant de démarrer le serveur.
// Le serveur a à disposition un fichier de configuration qui contient les adresses des serveurs.
// Il gère aussi un flag debug qui permet d'allonger le temps d'attente des requêtes avant de timeout.
package main

import (
	_ "embed"
	"flag"
	"log"
	"strconv"

	"github.com/Lazzzer/labo3-sdr/internal/server"
	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

//go:embed config.json
var config string
var DEBUG_DELAY = 3  // in seconds
var timeoutDelay = 1 // in seconds

// main est la méthode d'entrée du programme
func main() {

	debug := flag.Bool("debug", false, "Boolean: Run server in debug mode. Default is false")
	flag.Parse()
	if flag.Arg(0) == "" {
		log.Fatal("Invalid argument, usage: <server number>")
	}

	number, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatal("Invalid argument, usage: <server number>")
	}

	configuration, err := shared.Parse[types.Config](config)
	if err != nil {
		log.Fatal(err)
	}

	if number > len(configuration.Servers) || number < 0 {
		log.Fatal("Invalid server number")
	}

	if *debug {
		timeoutDelay *= DEBUG_DELAY
	}

	serv := server.Server{
		Debug:        *debug,
		DebugDelay:   DEBUG_DELAY,
		Number:       number,
		Address:      configuration.Servers[number],
		Servers:      configuration.Servers,
		TimeoutDelay: timeoutDelay,
	}
	serv.Run()
}
