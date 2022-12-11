package types

type Config struct {
	Address string         `json:"address,omitempty"` // Adresse du serveur
	Servers map[int]string `json:"servers"`           // Adresses des serveurs disponibles
}

type LogType string

const (
	INFO    LogType = "INFO"
	ERROR   LogType = "ERROR"
	MESSAGE LogType = "MESSAGE"
	COMMAND LogType = "COMMAND"
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
	Ann MessageType = "announcement" // Annonce
	Res MessageType = "result"       // Résultat d'une élection
)

type Process struct {
	Number int `json:"number"` // Numéro du processus
	Value  int `json:"value"`  // Valeur de la charge du processus
}

type Message struct {
	Type      MessageType `json:"message_type"` // Type du message
	Elected   int         `json:"elected"`      // Numéro du processus élu
	Processes []Process   `json:"processes"`    // Liste des processus
}
