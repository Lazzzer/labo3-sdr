package server

import (
	"fmt"
	"log"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func handleMessage(messageStr string) (string, error) {
	message, err := shared.Parse[types.Message](messageStr)
	if err != nil || message.Type == "" {
		return "", fmt.Errorf("invalid message")
	}

	// TODO : handle message for Chang & Roberts algorithm
	log.Println(message)

	return "message", nil
}
