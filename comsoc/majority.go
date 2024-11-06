package comsoc

func MajoritySWF(p Profile, _ []int) (count Count, err error) {

	count = make(Count)
	candid := make([]Alternative, len(p[0]))

	for i := 0; i < len(candid); i++ {
		candid[i] = p[0][i]
		count[Alternative(p[0][i])] = 0
	}

	for _, prefs := range p {
		// if err := CheckProfile(prefs, candid); err != nil {
		//	return nil, err
		// }
		count[prefs[0]]++
	}
	return count, nil
}

func MajoritySCF(p Profile, _ []int) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p, nil)
	if err != nil {
		return nil, err
	}
	return MaxCount(count), nil
}
