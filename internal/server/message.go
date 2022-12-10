package server

import (
	"fmt"

	"github.com/Lazzzer/labo3-sdr/internal/shared"
	"github.com/Lazzzer/labo3-sdr/internal/shared/types"
)

func handleMessage(messageStr string) (string, error) {
	message, err := shared.Parse[types.Message](messageStr)
	if err != nil || message.Type == "" {
		return "", fmt.Errorf("invalid message")
	}

	switch message.Type {
	case types.Ann:
		handleAnn(message)
	case types.Res:
		handleRes(message)
	}

	return "message", nil
}

func handleAnn(message *types.Message) {
	// TODO: handle ann message
}

func handleRes(message *types.Message) {
	// TODO: handle res message
}
