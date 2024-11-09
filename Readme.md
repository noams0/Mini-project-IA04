# Serveur de vote

|  Information |   Valeur              |
| :------------ | :------------- |
| **Auteurs** | Noam Seuret, Octave Leroy |



## Introduction

Ce projet implémente un serveur de vote en Go, avec plusieurs méthodes de vote et des options de configuration pour simuler des agents votants. Nous détaillons par la suite les méthodes de votent développer et comment installer ce serveur de votes. 

## I - Installation

Pour installer automatiquement ce serveur de vote, exécutez la commande suivante :

```bash
go get github.com/noams0/Mini-project-IA04
```

## II - Méthodes de votes

Le serveur prend en charge plusieurs méthodes de vote, accessibles en spécifiant le nom de la méthode lors de la configuration.

### 1 - Majorité

Le vote majoritaire est implémenté dans ce projet. Utilisez la règle `"majority"`pour activer cette méthode de vote.


### 2 - Borda

Le système de vote de Borda est disponible en utilisant la règle `"borda"`. Cette méthode classe les choix en attribuant des points à chaque option, selon leur classement par les votants. 

### 3 - Copeland 

La méthode de vote de Copeland peut être activée avec la règle `"copeland"`. Cette méthode compare chaque option en duel avec les autres pour établir un classement en fonction des victoires directes.

### 4 - Approval 

Le vote par approbation est disponible avec la règle `"approval"`. Dans cette méthode, chaque votant peut approuver plusieurs options. La première option spécifiée par chaque votant est utilisée comme seuil d'acceptation. Si aucune option ou une option incohérente est fournie, le vote de l’agent ne sera pas pris en compte.


### 5 - Condorcet

La méthode de Condorcet est accessible avec la règle `"condorcet"`. Cette méthode identifie un gagnant de Condorcet (le candidat qui bat tous les autres en duel) s’il existe, et fournit un classement des candidats basé sur leur proximité au gagnant de Condorcet. Si aucun gagnant de Condorcet n'existe, le classement de Copeland est renvoyé.



### 6 - STV

Le vote par système de vote transférable (STV) est disponible via la règle `"stv"`. Cette méthode élimine progressivement les candidats les moins populaires et redistribue leurs voix jusqu’à ce qu’un gagnant soit déterminé.

### 7 - Kemeny 

La méthode de consensus de Kemeny est disponible avec la règle `"kemeny"`. Cette méthode identifie le consensus optimal parmi les préférences des votants. Notez que cette méthode peut être longue à calculer lorsque le nombre d’alternatives est élevé.

Voici un complément pour ton fichier README en ajoutant des sections pour expliquer la structure du projet et les commandes :


## III - Structure du projet

Le projet est structuré en plusieurs dossiers principaux pour organiser les différentes fonctionnalités :

`comsoc` : Contient toutes les méthodes de vote implémentées, dont les règles de majorité, Borda, Copeland, Condorcet, STV, et Kemeny.

`agt` : Comprend le serveur principal ainsi que les différents agents (ballot, voter, admin). Cela gère la création des scrutins et le traitement des votes envoyés par les agents.

`cmd` : Ce dossier contient les commandes pour lancer les différentes parties de l’application.
- `server` : Permet de démarrer le serveur de vote.
- `launchAgent` : Utilisé pour initier les agents votants, enregistrer les votes, et récupérer les résultats une fois le scrutin terminé.

## IV - Démarrage du serveur et des agents

### Lancer le serveur

Pour démarrer le serveur, utilisez la commande suivante :
`go run cmd/server.go`

Le serveur sera lancé et écoutera les requêtes REST pour la création de scrutins, l’enregistrement de votes, et la récupération des résultats.

### Lancer les agents
Une fois le serveur lancé, vous pouvez soit effectué des véritable requêtes ou les lancer automatiquement avec la commande suivante :
`go run cmd/launchAgent.go`

Cette commande permet d’initier les agents, d'envoyer les votes au serveur et de récupérer les résultats une fois le vote terminé.


### Les requêtes 

### Commande `/new_ballot`


*[Extrait du sujet](https://gitlab.utc.fr/lagruesy/ia04/-/blob/main/docs/sujets/activit%C3%A9s/serveur-vote/api.md?ref_type=heads)*
- Requête : `POST`
- Objet `JSON` envoyé

| propriété  | type        | exemple de valeurs possibles                                      |
|------------|-------------|-------------------------------------------------------------------|
| `rule`      | `string`       | `"majority"`,`"borda"`, `"approval"`, `"stv"`, `"kemeny"`,... |
| `deadline`  | `string`       | `"2023-10-09T23:05:08+02:00"`  (format RFC 3339)                             |
| `voter-ids` | `[string,...]` | `["ag_id1", "ag_id2", "ag_id3"]`                                       |
| `#alts`     | `int`          | `12` |
| `tie-break` | `[int,...]`.   | `[4, 2, 3, 5, 9, 8, 7, 1, 6, 11, 12, 10]` |

- Code retour

| Code retour | Signification |
|-------------|---------------|
| `201`       | vote créé     |
| `400`       | bad request   |
| `501` 	  | not implemented |
*Exemple concret*

```Json
{
"rule": "majority",
"deadline": "2025-10-09T23:05:08+02:00",
"voter-ids": ["agt1", "agt2", "agt3"],
"#alts": 12,
"tie-break": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
}
```
Réponse `JSON`

Code HTTP : `201 Created`
```Json
{
"ballot-id": "scurtinNum0"
}
```
Au niveau du fichier launchAgent, c'est l'admin au travers de sa méthode StartVotingSession qui crée le vote, en faisant
appel à /new_ballot

```go
ballotID, _ := administrator.StartVotingSession("majority", deadline, voterIDs, 6, tb)
```

### Commande `/vote`

*[Extrait du sujet](https://gitlab.utc.fr/lagruesy/ia04/-/blob/main/docs/sujets/activit%C3%A9s/serveur-vote/api.md?ref_type=heads)*

- Requête : `POST`
- Objet `JSON` envoyé

| propriété   | type | exemple de valeurs possibles |
|------------|-------------|------------------------|
| `agent-id` | `string` | `"ag_id1"` |
| `ballot-id`| `string` | `"scrutin12"` |
| `prefs`    | `[int,...]` | `[1, 2, 4, 3]` |
| `options`  | `[int,...]` | `[3]` |

*Remarque :*`options`, dans notre cas, il ne sert uniquement qu'à la méthode de vote approval. `options` représente ainsi 
le seuil du nombre de candidats pour lequel l'agent votant donne son approval

- code retour

| Code retour | Signification |
|-------------|---------------|
| `200`       | vote pris en compte  |
| `400`       | bad request          |
| `403`       |	vote déjà effectué   |
| `501` 	  | Not Implemented      |
| `503`       | la deadline est dépassée |

*Exemple concret* 

```Json
{
"agent-id": "agt1",
"ballot-id": "scurtinNum0",
"prefs": [3, 4, 1, 5, 9, 8, 7, 2, 6, 11, 12, 10]
}
```

Réponse : Code HTTP : `200 OK`

Au niveau du fichier launchAgent, c'est un agent au travers de sa méthode vote qui fait son vote, en faisant
appel à /vote. Avec ballotID qui correspond à l'id du ballot pour qui il va voter

```go
ag.Vote(ballotID)
```

### Commande `/result`

*[Extrait du sujet](https://gitlab.utc.fr/lagruesy/ia04/-/blob/main/docs/sujets/activit%C3%A9s/serveur-vote/api.md?ref_type=heads)*

- Requête : `POST`
- Objet `JSON` envoyé

| propriété  | type | exemple de valeurs possibles                                  |
|------------|-------------|-----------------------------------------------------|
| `ballot-id`    | `string` | `"scrutin12"` |


- code retour

| Code retour | Signification   |
|-------------|-----------------|
| `200`       | OK              |
| `425`       | Too early       |
| `404`       |	Not Found       |

- Objet `JSON` renvoyé (si `200`)

| propriété   | type | exemple de valeurs possibles |
|------------|-------------|------------------------|
| `winner`   | `int`       | `4`                    |
| `ranking`  | `[int,...]` | `[2, 1, 4, 3]`         |

Au niveau du fichier launchAgent, on récupère cela grâce à l'administrator : 

```go
administrator.GetResults(ballot)
```