// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR
// source: https://twin.sh/articles/35/how-to-add-colors-to-your-console-terminal-output-in-go

package shared

import (
	"encoding/json"
	"log"

	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func Parse[T types.Config | types.Command | types.Message | types.Acknowledgement](jsonStr string) (*T, error) {
	var object T

	err := json.Unmarshal([]byte(jsonStr), &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

func Log(logType types.LogType, message string) {
	switch logType {
	case types.INFO:
		log.Println(CYAN + "(INFO) " + RESET + message)
	case types.ERROR:
		log.Println(RED + "(ERROR) " + RESET + message)
	case types.MESSAGE:
		log.Println(ORANGE + "(MESSAGE) " + RESET + message)
	case types.COMMAND:
		log.Println(YELLOW + "(COMMAND) " + RESET + message)
	}
}

// Variables pour colorer le texte dans la console
var RESET = "\033[0m"         // Variable pour réinitialiser la couleur du texte
var RED = "\033[31m"          // Variable pour colorer le texte en rouge
var PINK = "\033[38;5;198m"   // Variable pour colorer le texte en rose
var GREEN = "\033[32m"        // Variable pour colorer le texte en vert
var YELLOW = "\033[33m"       // Variable pour colorer le texte en jaune
var ORANGE = "\033[38;5;208m" // Variable pour colorer le texte en orange
var CYAN = "\033[36m"         // Variable pour colorer le texte en cyan
var BOLD = "\033[1m"          // Variable pour changer le texte en gras
