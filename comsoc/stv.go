package comsoc

import "fmt"

func StvSWF(p Profile, tiebreak []int) (c Count, err error) {

	fmt.Println("profil", p)

	// Inverse du tibreak pour avoir la pire possibilité
	tbinv := InverseAlternatives(TransformAlternatives(tiebreak))
	tb := TieBreakFactory(tbinv)

	// pour le SwF on peut faire un map avec le nombres de tour ou le candidat est en lice
	c = make(map[Alternative]int)

	// Alternatives du vote
	var alts []Alternative = make([]Alternative, len(p[0])) // Créer le slice des alternatives

	// mettre les alternatives a 0
	for i := 0; i < len(alts); i++ {
		alts[i] = Alternative(i + 1)
		c[Alternative(i+1)] = 0
	}

	// Créer le profile ou on enlevera les candidats au fur et à mesure
	var ptemp Profile = make(Profile, len(p))
	copy(ptemp, p)

	// On réalise un vote majoritaire
	// Si un candidat a la majorité des voix il est élu directement
	// Sinon on enlève le candidat avec le moins de voix

	for tour := 0; tour < len(alts)-1; tour++ {
		// n-1 tours, avec n le nombres de candidats

		// Récupère les résultats du vote majoritaire
		temp_count, _ := MajoritySWF(ptemp, nil)
		fmt.Println(temp_count)

		// Récupère l' ou les alternatives avec le plus de votes
		// best_alts, best_values := Maxvalue(temp_count)

		// Si elle a plus de la majorité
		// if best_values[0] > len(p2)/2 {
		// 	c[best_alts[0]] += 1 // ajoute un a cette alternative car elle gagne
		// 	return c, nil
		// }

		// récupère le ou les pires alternatives
		worstAlts := MinCount(temp_count)
		fmt.Println(worstAlts)

		// tie break sur les alternatives
		worstAlt, err := tb(worstAlts)
		if err != nil {
			return Count{}, err
		}

		// Elimines le premier des pires du profil
		for i, prefs := range ptemp {
			newPrefs := make([]Alternative, 0, len(prefs)-1)
			for _, alt := range prefs {
				if alt != worstAlt {
					// Enlève l'alternative des préférences de chaque votants
					// et le met dans p2
					newPrefs = append(newPrefs, alt)
				}
			}
			ptemp[i] = newPrefs

		}

		// Augmente de 1 toutes les alternatives non supprimés du count
		for i := range c {
			if i != worstAlt && c[i] == tour { // Si différents de la valeur qu'on élimine il faut ajouter 1 !
				c[i]++
			}
		}

	}

	fmt.Println("profil", p)
	return c, nil
}

func StvSCF(p Profile, tb []int) (bestAlts []Alternative, err error) {
	fmt.Println("Avant")
	fmt.Println("profil", p)
	c, ok := StvSWF(p, tb)
	fmt.Println("Après")
	err = ok

	if err != nil {
		return []Alternative{}, err
	}

	return MaxCount(c), err

}
