// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR
// source: https://twin.sh/articles/35/how-to-add-colors-to-your-console-terminal-output-in-go

// Package shared propose des fonctions utilitaires pour le projet.
package shared

import (
	"encoding/json"
	"log"
	"math"
	"strconv"

	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

// Parse permet de parser un objet JSON en un objet de type T.
func Parse[T types.Config | types.Command | types.Message | types.Acknowledgement](jsonStr string) (*T, error) {
	var object T

	err := json.Unmarshal([]byte(jsonStr), &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

// Log permet d'afficher un message dans la console avec une couleur différente selon le type de log.
func Log(logType types.LogType, message string) {
	switch logType {
	case types.INFO:
		log.Println(CYAN + "(INFO) " + RESET + message)
	case types.DEBUG:
		log.Println(ORANGE + "(DEBUG) " + RESET + message)
	case types.ERROR:
		log.Println(RED + "(ERROR) " + RESET + message)
	case types.MESSAGE:
		log.Println(PINK + "(MESSAGE C&R) " + RESET + message)
	case types.COMMAND:
		log.Println(YELLOW + "(COMMAND) " + RESET + message)
	}
}

// GetNbProcessWithMinValue retourne le numéro du processus avec la valeur la plus petite.
func GetNbProcessWithMinValue(processes *[]types.Process) int {
	minValue := math.MaxInt
	minProcessNumber := -1

	for _, p := range *processes {
		if p.Value < minValue {
			minValue = p.Value
			minProcessNumber = p.Number
		}
	}

	return minProcessNumber
}

// ShowProcessList retourne une chaîne de caractères représentant la liste des processus.
func ShowProcessList(processes []types.Process, withValue bool) string {
	var list string
	list = "["
	for i, p := range processes {
		list += "P" + strconv.Itoa(p.Number)
		if withValue {
			list += ":" + strconv.Itoa(p.Value)
		}
		if i != len(processes)-1 {
			list += ", "
		}
	}
	list += "]"
	return list
}

// Variables pour colorer le texte dans la console
var RESET = "\033[0m"         // Variable pour réinitialiser la couleur du texte
var RED = "\033[31m"          // Variable pour colorer le texte en rouge
var PINK = "\033[38;5;219m"   // Variable pour colorer le texte en rose
var PURPLE = "\033[38;5;198m" // Variable pour colorer le texte en violet

var GREEN = "\033[32m"        // Variable pour colorer le texte en vert
var YELLOW = "\033[33m"       // Variable pour colorer le texte en jaune
var ORANGE = "\033[38;5;208m" // Variable pour colorer le texte en orange
var CYAN = "\033[36m"         // Variable pour colorer le texte en cyan
var BOLD = "\033[1m"          // Variable pour changer le texte en gras
