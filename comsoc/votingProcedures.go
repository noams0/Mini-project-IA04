package comsoc

import (
	"sort"
)

func SWF(p Profile) (count Count, err error) {
	count = make(Count)                      // Initialiser le décompte
	candid := make([]Alternative, len(p[0])) // Alternatives candidates

	for _, prefs := range p {
		if err := CheckProfile(prefs, candid); err != nil {
			return nil, err
		}
		for _, alt := range prefs {
			count[alt]++
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

func SWFFactory(swf func(p Profile, t []int) (Count, error), tieBreak func([]Alternative) (Alternative, error)) func(Profile, []int) ([]Alternative, error) {

	return func(p Profile, thresholds []int) ([]Alternative, error) {

		if thresholds == nil { // pas de thresholds dans cette fonction
			thresholds = []int{} // slice vide par défaut
		}

		count, err := swf(p, thresholds)

		if err != nil {
			return nil, err
		}
		alts := make([]Alternative, 0, len(count))

		// Ajoutez les clés de count dans le slice
		for alt := range count {
			alts = append(alts, alt)
		}

		// Tri par ordre décroissant de valeurs de count, avec résolution des égalités
		sort.SliceStable(alts, func(i, j int) bool {
			if count[alts[i]] != count[alts[j]] {
				return count[alts[i]] > count[alts[j]]
			}
			// Cas d'égalité : on utilise tieBreak pour départager
			winner, err := tieBreak([]Alternative{alts[i], alts[j]})
			if err != nil {
				return false // En cas d'erreur, on peut décider de ne rien changer
			}
			return winner == alts[i]
		})
		return alts, nil
	}
}

func SCFFactory(scf func(p Profile, t []int) ([]Alternative, error), tb func([]Alternative) (Alternative, error)) func(Profile, []int) (Alternative, error) {
	return func(p Profile, thresholds []int) (Alternative, error) {

		if thresholds == nil {
			thresholds = []int{}
		}

		res, err := scf(p, thresholds)
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
