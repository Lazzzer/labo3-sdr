# labo3-sdr

## Auteurs

Lazar Pavicevic et Jonathan Friedli

## Contexte

Ce projet est réalisé dans le cadre du cours de Systèmes Distribués et Répartis (SDR) de la HEIG-VD.

Dans ce laboratoire, nous allons implémenter l'algorithme de Chang et Roberts avec panne de processus afin de déterminer le processus de meilleur aptitude  dans un groupe de processus. Toutes les connexions seront réalisées en UDP.

## Client

### Commandes

```bash
# Commande ajoutant une charge (valeur entière) à un serveur
charge <idServer> <value>

# Commande demandant à un serveur quel est le processus élu
whoTheBoss <idServer>

# Commande demandant une élection à un serveur
election <idServer>

# Commande permettant de simuler la panne d'un server
kill <idServer>

```