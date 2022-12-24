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

# Les tests

Les tests peuvent être lancés avec les commandes suivantes:

```bash
# A la racine du projet
go test -race ./test/. -v

# Si besoins, en vidant le cache
go clean -testcache && go test -race ./test/. -v
```

Notre fichier de test comporte un `TestClient` et 3 serveurs de test (dont 1 down) composant le réseau, lancés dans la fonction `init()`. Le TestClient peut envoyer des commandes à un serveur de test et vérifier la string de réponse avec un résultat attendu.

Seul un test est lancé, `TestElectionWithDownServers`, qui va faire dans l'ordre:

- Ajouter une charge au serveur 1 et 3. A la fin de cette étape, le serveur 1 a une charge plus élevée que le serveur 3.
- Demander une nouvelle élection au serveur 1
- Demander au serveur 1 quel est le processus élu de la dernière élection. Comme il y a une élection en cours, le serveur répondra à la fin de l'élection. Le serveur répondra que le serveur 3 (avec le Processus P2) est élu.

Ce test vérifie alors que l'élection peut se faire même en cas de panne avant une élection.

![Tests](/docs/tests.png)

## Procédure de tests manuels

TODO

## Implémentation

### Le client

Le client effectue une nouvelle connexion UDP à un serveur à chaque commande envoyée. Pour chaque commande, il attend une réponse du serveur et affiche le résultat dans la console. En cas de timeout, il notifie l'utilisateur que le serveur est down.

Le client parse l'input en ligne de commande et crée un objet `Command` si l'input est valide. Il transforme ensuite cet objet en string JSON et l'envoie au serveur.

Quitter un client avec CTRL+C ou en envoyant la commande `quit` ferme la connexion UDP en cours et arrête le client gracieusement.

### Le serveur

Au niveau des spécificités du serveur, ce dernier répond à chaque communication reçue, que ce soit par un `Acknowledgement` pour un `Message` inter-server, ou une simple string pour un `Command`.

A part dans le cas d'une commande `ask`, le serveur envoie juste une réponse générique spécifiant que la commande a bien été reçue. Ainsi, nous pouvons mettre en place un timeout côté client pour savoir si le serveur est down.

Lorsqu'un serveur ne répond pas à un `Message` inter-server, il est considéré comme down et l’émetteur envoie le message au serveur suivant dans la liste des serveurs.

Ce système de timeout permet de gérer certains cas de panne de processus mais pas tous. En effet, il se peut qu'une élection puisse se bloquer lors de cas où une réception de message (que ce soit une annonce ou un résultat) a pu se faire (c'est-à-dire, avec envoi d'ack) juste avant que ce dernier ne tombe en panne et ne puisse transmettre l'information au prochain.

### Points à améliorer

Comme brièvement évoqué dans la partie serveur, il y a des scénarios où une panne peut s'avérer problématique lors d'une élection. Nous pourrions compléter l'algorithme d'élection pour qu'il puisse gérer ce genre de cas.

Autre point, nous utilisons deux fichiers `config.json` pour les serveurs et les clients qui sont identiques. Cela est principalement dû aux limitations imposées par `go:embed` qui nous permet de build le fichier dans l'exécutable. Nous pourrions alors soit fusionner les deux fichiers en remaniant notre approche avec le package `embed` soit rajouter des configurations spécifiques pour chaque élément afin de justifier des fichiers séparés.
