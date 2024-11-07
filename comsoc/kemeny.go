package comsoc

import (
	"fmt"
	"gitlab.utc.fr/lagruesy/ia04/utils"
	"sort"
)

//Rajouter gestion d'erreurs

func DistanceEdition(pref1, pref2 []Alternative) int {
	distance := 0
	n := len(pref1)

	indexMap1 := make(map[Alternative]int)
	for i, v := range pref1 {
		indexMap1[v] = i
	}
	indexMap2 := make(map[Alternative]int)
	for i, v := range pref2 {
		indexMap2[v] = i
	}
	var a Alternative = 0
	var b Alternative = 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a, b = pref1[i], pref1[j]
			if (indexMap1[a] > indexMap1[b] && indexMap2[a] < indexMap2[b]) || (indexMap1[a] < indexMap1[b] && indexMap2[a] > indexMap2[b]) {
				distance += 1
			}
		}
	}
	return distance
}

func DistanceEditionProfile(profile Profile, rangement []Alternative) int {
	total_distance := 0
	dist := 0
	for _, pref := range profile {
		dist = DistanceEdition(pref, rangement)
		total_distance += dist
	}
	return total_distance
}

func indexOf(alternative int, ranking []int) int {
	for i, alt := range ranking {
		if alt == alternative {
			return i
		}
	}
	return len(ranking)
}

func KemenySWF(profile Profile, tb []int) (c Count, err error) {
	var n int = len(profile[0])
	perm := utils.FirstPermutation(n)
	for i, _ := range perm {
		perm[perm[i]] += 1
	}
	dist := DistanceEditionProfile(profile, IntSliceToAlternativeSlice(perm))
	distMinimal := dist

	rangementsCons := [][]int{}
	rangementsCons = append(rangementsCons, perm)
	perm, ok := utils.NextPermutation(perm)
	for ok {

		dist = DistanceEditionProfile(profile, IntSliceToAlternativeSlice(perm))
		if distMinimal > dist {
			distMinimal = dist
			rangementsCons = [][]int{}
			rangementsCons = append(rangementsCons, perm)
		}
		if distMinimal == dist {
			rangementsCons = append(rangementsCons, perm)
		}
		perm, ok = utils.NextPermutation(perm)
	}
	consensus := []int{}
	if len(rangementsCons) == 0 {
		consensus = rangementsCons[0]
	} else {
		fmt.Println("Avant T-B : ", rangementsCons)

		// Trier rangementsCons selon l'ordre de tie-break `tb`
		sort.SliceStable(rangementsCons, func(i, j int) bool {
			for _, alternative := range tb {
				posI := indexOf(alternative, rangementsCons[i])
				posJ := indexOf(alternative, rangementsCons[j])

				if posI != posJ {
					return posI < posJ
				}
			}
			return false
		})

		// Choisir le premier élément après tri par tie-break
		consensus = rangementsCons[0]
	}
	fmt.Println(consensus)
	count := make(Count)
	for i := 0; i < len(consensus); i++ {
		count[Alternative(consensus[i])] = len(consensus) - i
	}
	fmt.Println(count)
	return count, nil
}

func KemenySCF(p Profile, tb []int) (bestAlts []Alternative, err error) {
	c, ok := KemenySWF(p, tb)
	err = ok
	if err != nil {
		return []Alternative{}, err
	}

	return MaxCount(c), err

}
