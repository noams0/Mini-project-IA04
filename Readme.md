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
