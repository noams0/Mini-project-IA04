package comsoc

import (
	"errors"
	"fmt"
)

func TieBreak(alts []Alternative) (Alternative, error) {
	if len(alts) == 0 {
		return 0, fmt.Errorf("no alternatives provided")
	}
	// Renvoie la première alternative
	return alts[0], nil
}

func TieBreakFactory(orderedAlts []Alternative) func([]Alternative) (Alternative, error) {
	return func(candidates []Alternative) (Alternative, error) {
		if len(candidates) == 0 {
			return 0, errors.New("la liste des candidats est vide")
		}
		// map utilisé pour trouver la position d'un candidat dans l'ordre prédéfini
		ranking := make(map[Alternative]int)
		for i, alt := range orderedAlts {
			ranking[alt] = i
		}

		// trouver le candidat avec le rang le plus bas
		winner := candidates[0]
		minRank := ranking[winner]
		for _, candidate := range candidates {
			if rank, found := ranking[candidate]; found {
				if rank < minRank {
					winner = candidate
					minRank = rank
				}
			} else {
				return 0, fmt.Errorf("le candidat %v n'est pas dans la liste de séquence prédéfinie", candidate)
			}
		}

		return winner, nil
	}
}
