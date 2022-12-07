package types

type Config struct {
	Address string         `json:"address,omitempty"` // Adresse du serveur
	Servers map[int]string `json:"servers"`           // Adresses des serveurs disponibles
}
