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

### Test n° 1

Durant ce test, nous allons simplement vérifier que l'élection se passe comme prévu. Nous pouvons faire ce test en mode `debug` ou pas.

Pour ce faire, nous lançons 3 serveurs et 1 client. Nous ajoutons ensuite des charges sur les serveurs. Une charge de 30 sur le premier, 15 sur le 2ème et 35 sur le 3ème. Nous demandons ensuite une nouvelle élection au serveur 1 puis nous demandons au serveur 3 quel est le processus élu. Nous devrions obtenir le processus P2. En effet, ce dernier a la plus petite charge.

```bash
# Input du client
1 add 30
2 add 15
3 add 35
1 new
3 ask # après le premier tour ou à la fin de l'élection
```

Résultat obtenu:

![Test 1](/docs/test1.png)

### Test n° 2

Durant ce test, nous allons vérifier que l'élection se passe bien même si un serveur est down. Nous pouvons faire ce test en mode `debug` ou pas.

Pour ce faire, nous lançons 3 serveurs et 1 client. Nous ajoutons ensuite des charges sur les serveurs. Une charge de 30 sur le premier, 15 sur le 2ème et 35 sur le 3ème. Logiquement le serveur 2 devrait être élu. Cependant nous allons simuler une panne sur ce serveur puis demander une élection au serveur 1 et finir par lui demander le processus élu. Nous devrions obtenir le processus P1. En effet, ce dernier a la plus petite charge après la panne du serveur 2.

```bash
# Input du client
1 add 30
2 add 15
3 add 35
2 stop
1 new
1 ask # après le premier tour ou à la fin de l'élection
```

Résultat obtenu:

![Test 2](/docs/test2.png)

### Test n° 3

Nous allons maintenant lancer les serveurs et le client en mode `debug` afin de les ralentir et de simuler un temps de traitement assez long. Nous lançons donc 3 serveurs et 2 clients. Nous ajoutons donc les mêmes charge qu'auparavant. Depuis le client 1, nous demandons une nouvelle élection au serveur 1. Une fois que le premier tour de serveur est terminé, nous demandons depuis le client 2 une nouvelle élection au serveur 2. La deuxième élection va juste être annulée car le serveur 1 a déjà lancé une élection qui a fini le premier tour.

```bash
# Input du client1
1 add 30
2 add 15
3 add 35
1 new
```

```bash
# Input du client2
2 new # après le premier tour de l'élection
```

Résultat obtenu:

![Test 3](/docs/test3.png)

### Test n° 4

Nous allons lancer la même configuration que durant le test n° 3 mais cette fois-ci, nous allons lancer les deux élections de manière simultanée. Les deux élections vont donc commencer par tourner en parallèle puis vu qu'elles vont élire le même processus (P1), les élections vont se terminer en même temps avec le même résultat.

```bash
# Input du client1
1 add 30
2 add 15
3 add 35
1 new # Simultanément avec le client 2
```

```bash
# Input du client2
2 new # Simultanément avec le client 1
```

Résultat obtenu:

![Test 4](/docs/test4.png)

### Test n° 5

Le test n° 5 est une variante du 4 où nous ne rajoutons pas de charges aux serveurs. Ici, les deux élections vont donc commencer par tourner en parallèle et élire un processus différent. Pour corriger ce problème, l'élection qui a le numéro du processus élu le plus petit va être répétée et va resynchroniser le processus élu pour tous les serveurs du réseau.

```bash
# Input du client1
1 new # Simultanément avec le client 2
```

```bash
# Input du client2
2 new # Simultanément avec le client 1
```

Résultat obtenu:

![Test 5](/docs/test5.png)

### Test n° 6

Nous allons lancer la même configuration que durant le test n° 3. Le client 1 va lancer une élection et juste après, le client 2 va ajouter une grande charge sur le serveur 2. De cette manière le serveur 2 aura la plus grande charge. Cependant, au démarrage de l'élection, le serveur 2 était celui avec la charge la plus faible. Il va donc être élu.

```bash
# Input du client1
1 add 30
2 add 15
3 add 35
1 new
```

```bash
# Input du client2
2 add 50 # Juste après la demande d'élection du client 1
```

Résultat obtenu:

![Test 6](/docs/test6.png)

## Implémentation

### Le client

Le client effectue une nouvelle connexion UDP à un serveur à chaque commande envoyée. Pour chaque commande, il attend une réponse du serveur et affiche le résultat dans la console. En cas de timeout, il notifie l'utilisateur que le serveur est down.

Le client parse l'input en ligne de commande et crée un objet `Command` si l'input est valide. Il transforme ensuite cet objet en string JSON et l'envoie au serveur.

Quitter un client avec CTRL+C ou en envoyant la commande `quit` ferme la connexion UDP en cours et arrête le client gracieusement.

### Le serveur

Au niveau des spécificités du serveur, ce dernier répond à chaque communication reçue, que ce soit par un `Acknowledgement` pour un `Message` inter-server, ou une simple string pour un `Command`.

A part dans le cas d'une commande `ask`, le serveur envoie juste une réponse générique spécifiant que la commande a bien été reçue. Ainsi, nous pouvons mettre en place un timeout côté client pour savoir si le serveur est down.

Lorsqu'un serveur ne répond pas à un `Message` inter-server, il est considéré comme down et l’émetteur envoie le message au serveur suivant dans la liste des serveurs.

Ce système de timeout permet de gérer certains cas de panne de processus mais pas tous. En effet, il se peut qu'une élection puisse se bloquer lors de cas où une réception de message (que ce soit une annonce ou un résultat) a pu se faire (c'est-à-dire, avec envoi d'ack) juste avant que le serveur ne tombe en panne et ne puisse transmettre l'information au prochain.

### Points à améliorer

Comme brièvement évoqué dans la partie serveur, il y a des scénarios où une panne peut s'avérer problématique lors d'une élection. Nous pourrions compléter l'algorithme d'élection pour qu'il puisse gérer ce genre de cas.

Autre point, nous utilisons deux fichiers `config.json` pour les serveurs et les clients qui sont identiques. Cela est principalement dû aux limitations imposées par `go:embed` qui nous permet de build le fichier dans l'exécutable. Nous pourrions alors soit fusionner les deux fichiers en remaniant notre approche avec le package `embed` soit rajouter des configurations spécifiques pour chaque élément afin de justifier des fichiers séparés.
