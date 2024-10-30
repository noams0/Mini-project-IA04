package comsoc

import "log"

func MajoritySWF(p Profile, _ []int) (count Count, err error) {
	count = make(Count)
	candid := make([]Alternative, len(p[0]))
	for i := range p[0] {
		count[Alternative(i+1)] = 0
	}
	for _, prefs := range p {
		if err := CheckProfile(prefs, candid); err != nil {
			return nil, err
		}
		count[prefs[0]]++
	}
	log.Println(count)
	return count, nil
}

func MajoritySCF(p Profile, _ []int) (bestAlts []Alternative, err error) {
	count, err := MajoritySWF(p, nil)
	if err != nil {
		return nil, err
	}
	return MaxCount(count), nil
}
