// Auteurs: Jonathan Friedli, Lazar Pavicevic
// Labo 3 SDR

// Package types propose différents types utilisés par l'application pour parser le fichier de configuration, les messages et les commandes.
package types

// Config représente la configuration du réseau de serveurs.
type Config struct {
	Address string         `json:"address,omitempty"` // Adresse du serveur
	Servers map[int]string `json:"servers"`           // Adresses des serveurs disponibles
}

type LogType string // Type de log

const (
	INFO    LogType = "INFO"    // Log d'information
	DEBUG   LogType = "DEBUG"   // Log de debug
	ERROR   LogType = "ERROR"   // Log d'erreur
	MESSAGE LogType = "MESSAGE" // Log de message
	COMMAND LogType = "COMMAND" // Log de commande
)

type CommandType string // Type de commande

const (
	Add  CommandType = "add"  // Commande d'ajout de valeur
	Ask  CommandType = "ask"  // Commande de demande du processus élu
	New  CommandType = "new"  // Commande de nouvelle élection
	Stop CommandType = "stop" // Commande d'arrêt du serveur
	Quit CommandType = "quit" // Commande de fermeture du client
)

// Command représente une commande envoyée par un client.
type Command struct {
	Type  CommandType `json:"command_type"`    // Type de la commande
	Value *int        `json:"value,omitempty"` // Valeur à ajouter
}

type MessageType string // Type de message

const (
	Ann MessageType = "announcement" // Message d'annonce
	Res MessageType = "result"       // Message de résultat d'élection
)

// Process représente un processus dans le réseau de serveurs.
type Process struct {
	Number int `json:"number"` // Numéro du processus
	Value  int `json:"value"`  // Valeur de la charge du processus
}

// Message représente un message envoyé par un serveur à un autre serveur.
type Message struct {
	Type      MessageType `json:"message_type"` // Type du message
	Elected   int         `json:"elected"`      // Numéro du processus élu
	Processes []Process   `json:"processes"`    // Liste des processus
}

// Acknowledgement représente un message d'acknowledgement envoyé par un serveur à un autre serveur pour confirmer la réception d'un message.
type Acknowledgement struct {
	From int `json:"number"` // Numéro du processus qui a envoyé l'ack
}
