// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

package shared

import (
	"encoding/json"
	"log"

	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func ParseConfig(configStr string) *types.Config {
	var config types.Config

	err := json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
