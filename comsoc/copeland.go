package comsoc

func CopelandSWF(p Profile) (Count, error) {
	candid := p[0]
	count := make(Count)

	for i := 0; i < len(candid)-1; i++ {
		for j := i + 1; j < len(candid); j++ {
			c1 := candid[i]
			c2 := candid[j]
			if compareTwoAlternatives(p, c1, c2) {
				count[c1]++
			} else {
				count[c2]++
			}
		}
	}

	return count, nil
}

func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := CopelandSWF(p)
	if err != nil {
		return nil, err
	}
	return MaxCount(count), nil
}
