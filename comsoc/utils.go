package comsoc

import (
	"errors"
	"fmt"
	"math"
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

// Pour STV

func MinCount(count Count) (bestAlts []Alternative) {
	// Trouve les min du Count
	var min int = math.MaxInt64

	for alt, value := range count {
		if value < min {
			min = value
			bestAlts = []Alternative{}       // Remet à 0 les meilleurs alternatives
			bestAlts = append(bestAlts, alt) // remet dans l'objet car peut changer d'emplacement mémoire

		} else if value == min {
			bestAlts = append(bestAlts, Alternative(alt))
		}
		// Sinon il se passe rien
	}
	return
}

func Maxvalue(count Count) (bestAlts []Alternative, bestValues []int) {
	// Trouve les max du Count

	var max int = -1
	for alt, value := range count {
		if value > max {
			max = value
			bestAlts = []Alternative{}       // Remet à 0 les meilleurs alternatives
			bestAlts = append(bestAlts, alt) // remet dans l'objet car peut changer d'emplacement mémoire
			bestValues = []int{}
			bestValues = append(bestValues, value)

		} else if value == max {
			bestAlts = append(bestAlts, alt)
			bestValues = append(bestValues, value)
		}
		// Sinon il se passe rien
	}
	return
}

func TransformInt(tiebreak []Alternative) (c []int) {
	c = make([]int, len(tiebreak))

	for i, v := range tiebreak {
		c[i] = int(v)
	}
	return
}

func InverseAlternatives(tiebreak []Alternative) (c []Alternative) {
	c = make([]Alternative, len(tiebreak))

	for i, v := range tiebreak {
		c[len(tiebreak)-i-1] = Alternative(v)
	}
	return
}

func IntSliceToAlternativeSlice(intSlice []int) []Alternative {
	alternativeSlice := make([]Alternative, len(intSlice))
	for i, v := range intSlice {
		alternativeSlice[i] = Alternative(v)
	}
	return alternativeSlice
}
