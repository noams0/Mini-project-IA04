package comsoc

func BordaSWF(p Profile) (count Count, err error) {
	count = make(Count)
	candid := make([]Alternative, len(p[0]))
	for _, prefs := range p {
		if err := CheckProfile(prefs, candid); err != nil {
			return nil, err
		}
		for i := 0; i < len(prefs); i++ {
			count[prefs[i]] += len(prefs) - i - 1
		}
	}
	return count, nil
}

func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
	count, err := BordaSWF(p)
	if err != nil {
		return nil, err
	}
	return MaxCount(count), nil
}
