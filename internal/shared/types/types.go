package types

type Config struct {
	Address string         `json:"address,omitempty"` // Adresse du serveur
	Servers map[int]string `json:"servers"`           // Adresses des serveurs disponibles
}

type LogType string

const (
	INFO  LogType = "INFO"
	ERROR LogType = "ERROR"
)

type CommandType string

const (
	Add  CommandType = "add"
	Ask  CommandType = "ask"
	New  CommandType = "new"
	Stop CommandType = "stop"
	Quit CommandType = "quit"
)

type Command struct {
	Type  CommandType `json:"command_type"`    // Type de la commande
	Value *int        `json:"value,omitempty"` // Valeur à ajouter
}

type MessageType string

const (
	REQ MessageType = "REQ" // Requête
)

type Message struct {
	Type MessageType `json:"message_type"` // Type du message
	From int         `json:"from"`         // Numéro du serveur qui a envoyé le message
}
