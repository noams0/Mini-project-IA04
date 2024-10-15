package comsoc

func MajoritySWF(p Profile) (count Count, err error) {
	count = make(Count)
	candid := make([]Alternative, len(p[0]))
	for _, prefs := range p {
		if err := CheckProfile(prefs, candid); err != nil {
			return nil, err
		}
		count[prefs[0]]++
	}
	return count, nil
}

func MajoritySCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p)
	if err != nil {
		return nil, err
	}
	return MaxCount(count), nil
}
