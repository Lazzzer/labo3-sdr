# Laboratoire 3 de SDR - Application de l'algorithme de Chang et Roberts

## Auteurs

Lazar Pavicevic et Jonathan Friedli

## Contexte

Ce projet est réalisé dans le cadre du cours de Systèmes Distribués et Répartis (SDR) de la HEIG-VD.

Dans ce laboratoire, nous implémentons l'algorithme de Chang et Roberts avec panne de processus afin de déterminer le processus avec la charge la moins élevée dans un groupe de processus. Toutes les connexions sont réalisées en UDP.

## Utilisation du programme

L'application contient deux exécutables : un pour le serveur et un pour le client.

Le client et le serveur peuvent être lancés en mode `debug`.

Le mode `debug` ralentit artificiellement d'un nombre paramétrable de secondes le serveur lorsqu'il traite et envoie des messages. Pour le client, l'activation du mode prolonge la durée d'un timeout avant de considérer un serveur comme down.

### Pour lancer un serveur:

Le serveur a besoin d'un entier en argument qui représente la clé des maps présentes dans son fichier de configuration. Ces maps indiquent l'adresse de tout les autres serveurs composant le réseau.

Il dispose du flag `--debug`.

```bash
# A la racine du projet

# Lancement du serveur n°1
go run cmd/server/main.go 1

# Lancement du serveur n°1 en mode race & debug
go run -race cmd/server/main.go --debug 1
```

### Pour lancer un client:

Le client n'a pas besoins d'arguments pour être lancé. Il peut cependant prendre un flag `--debug` pour augmenter la durée du timeout avant de considérer un serveur comme down pour notamment pouvoir tester le comportement du réseau quand du délai est ajouté.

```bash
# A la racine du projet

# Lancement d'un client
go run cmd/client/main.go

# Lancement d'un client en mode race & debug
go run -race cmd/client/main.go --debug
```

### Usages:

```bash
# A la racine du projet
go run cmd/server/main.go --help
go run cmd/client/main.go --help

# Ou si le projet a été compilé et que l'exécutable se trouve dans le dossier courant
.\main.exe --help # Sous Windows
./main --help # Sous Linux/macOS
```

Résultat pour le serveur:

```bash
Usage of ./main:
  -debug
    	Boolean: Run server in debug mode. Default is false
```

Résultat pour le client:

```bash
Usage of ./main:
  -debug
    	Boolean: Run client in debug mode. Default is false
```

### Commandes disponibles:

```bash

# Commande ajoutant une charge (valeur entière) à un serveur
<server number> add <number>

# Commande demandant à un serveur quel est le processus élu de la dernière élection (attend la réponse du serveur)
<server number> ask

# Commande demandant une nouvelle élection à un serveur
<server number> new

# Commande demandant à un serveur de s'arrêter, simulant une panne
<server number> stop

# Commande permettant de quitter le client
quit
```
