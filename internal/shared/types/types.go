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

type CommandType string // TODO: Maybe rename to ClientCommandType if we add server commands?

const (
	Add  CommandType = "add"
	Ask  CommandType = "ask"
	New  CommandType = "new"
	Stop CommandType = "stop"
)

type Command struct {
	Type   CommandType `json:"type"`            // Type de la commande
	Server int         `json:"server"`          // Numéro du serveur
	Value  *int        `json:"value,omitempty"` // Valeur à ajouter
}
