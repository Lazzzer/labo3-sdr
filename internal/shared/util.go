// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

package shared

import (
	"encoding/json"

	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func Parse[T types.Config | types.Command | types.Message](jsonStr string) (*T, error) {
	var object T

	err := json.Unmarshal([]byte(jsonStr), &object)
	if err != nil {
		return nil, err
	}

	return &object, nil
}
