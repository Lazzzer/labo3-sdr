// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

package shared

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"go/types"
)

func Parse(content string) *types.Config {
	var config types.Config

	err := json.Unmarshal([]byte(content), &config)

	if err != nil {
		fmt.Println(err)
		panic("Error: Could not parse object")
	}

	return &config
}
