package comsoc

import (
	"errors"
	"fmt"
	"math/rand"
)

func GeneratePreferences(numAgents, numAlternatives int) [][]Alternative {
	profile := make([][]Alternative, numAgents)
	for i := 0; i < numAgents; i++ {
		prefs := make([]Alternative, numAlternatives)

		for j := 0; j < numAlternatives; j++ {
			prefs[j] = Alternative(j + 1)
		}
		for j := range prefs {
			k := rand.Intn(numAlternatives)
			prefs[j], prefs[k] = prefs[k], prefs[j]
		}
		profile[i] = prefs
	}
	return profile
}

// renvoie l'indice ou se trouve alt dans prefs
func Rank(alt Alternative, prefs []Alternative) int {
	for i, a := range prefs {
		if alt == a {
			return i
		}
	}
	return -1
}

// renvoie vrai ssi alt1 est préférée à alt2
func IsPref(alt1, alt2 Alternative, prefs []Alternative) bool {
	for _, a := range prefs {
		if alt1 == a {
			return true
		}
		if alt2 == a {
			return false
		}
	}
	return false
}

// renvoie les meilleures alternatives pour un décompte donné
func MaxCount(count Count) (bestAlts []Alternative) {
	var maxVal = -1
	for _, val := range count {
		if val > maxVal {
			maxVal = val
		}
	}

	for alt, val := range count {
		if val == maxVal {
			bestAlts = append(bestAlts, alt)
		}
	}
	return
}

// vérifie les préférences d'un agent, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois
func CheckProfile(prefs []Alternative, alts []Alternative) error {
	if len(prefs) != len(alts) {
		return errors.New("les préférences ne sont pas complètes")
	}

	seen := make([]bool, len(alts))
	for _, pref := range prefs {
		if pref < 0 || pref > Alternative(len(alts)) {
			return fmt.Errorf("alternative %d n'est pas valide", pref)
		}
		if seen[pref-1] {
			return fmt.Errorf("alternative %d apparaît plus d'une fois", pref)
		}
		seen[pref-1] = true
	}

	return nil
}

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func CheckProfileAlternative(prefs Profile, alts []Alternative) error {
	for _, agentPrefs := range prefs {
		err := CheckProfile(agentPrefs, alts)
		if err != nil {
			return errors.New("erreur dans le profil de l'agent " + ": " + err.Error())
		}
	}
	return nil
}

func compareTwoAlternatives(p Profile, alt1 Alternative, alt2 Alternative) bool {
	var alt1Wins = 0
	var alt2Wins = 0
	for _, prefs := range p {
		for _, pref := range prefs {
			if pref == alt1 {
				alt1Wins++
				break
			}
			if pref == alt2 {
				alt2Wins++
				break
			}
		}
	}
	return alt1Wins > alt2Wins
}
