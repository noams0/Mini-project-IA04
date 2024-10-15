package comsoc

func CondorcetWinner(p Profile) (bestAlts []Alternative, err error) {
	candid := p[0]
	eliminated := make(map[Alternative]bool) // Map pour suivre les candidats éliminés
	for _, c := range candid {
		// Si le candidat a déjà été éliminé, on passe au suivant
		if eliminated[c] {
			continue
		}
		isWinner := true
		for _, c2 := range candid {
			if c != c2 {
				if !compareTwoAlternatives(p, c, c2) {
					isWinner = false
					break // Si c est battu par c2, on peut arrêter la comparaison
				} else {
					eliminated[c2] = true // Sinon continue, et en plus on sait que C2 n'est pas un gagnant de condorcet
				}
			}
		}
		if isWinner {
			return []Alternative{c}, nil
		}
	}
	return nil, nil
}
