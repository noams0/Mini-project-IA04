package comsoc

func SWF(p Profile) (count Count, err error) {
	count = make(Count)                      // Initialiser le décompte
	candid := make([]Alternative, len(p[0])) // Alternatives candidates

	// Vérifier le profil de chaque votant
	for _, prefs := range p {
		if err := CheckProfile(prefs, candid); err != nil {
			return nil, err
		}
		// Compter chaque vote dans le profil
		for _, alt := range prefs {
			count[alt]++ // Incrémenter le compteur pour chaque alternative votée
		}
	}
	return count, nil
}

func SCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := SWF(p)
	if err != nil {
		return nil, err
	}
	return MaxCount(count), nil
}

func SWFFactory(swf func(p Profile) (Count, error), tieBreak func([]Alternative) (Alternative, error)) func(Profile) ([]Alternative, error) {
	return func(p Profile) ([]Alternative, error) {
		count, err := swf(p)
		if err != nil {
			return nil, err
		}

		maxCount := MaxCount(count)
		if len(maxCount) == 1 {
			return maxCount, nil
		}

		winner, err := tieBreak(maxCount)
		if err != nil {
			return nil, err
		}
		return []Alternative{winner}, nil
	}
}

func SCFFactory(scf func(p Profile) ([]Alternative, error), tb func([]Alternative) (Alternative, error)) func(Profile) (Alternative, error) {
	return func(p Profile) (Alternative, error) {
		res, err := scf(p)
		if err != nil {
			return -1, err
		}
		winner, err := tb(res)
		if err != nil {
			return -1, err
		}
		return winner, nil
	}
}
