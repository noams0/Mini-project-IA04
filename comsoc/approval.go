package comsoc

import (
	"fmt"
)

func ApprovalSWF(p Profile, thresholds []int) (count Count, err error) {
	count = make(Count)

	if len(thresholds) != len(p) {
		return nil, fmt.Errorf("not the same number")
	}

	fmt.Println(len(p[0]))
	for i := range p[0] {
		count[Alternative(i+1)] = 0
	}
	fmt.Println(p)

	for ind, prefs := range p {

		if err := CheckProfile(prefs, p[0]); err != nil {
			return nil, err
		}
		if thresholds[ind] > len(prefs) {
			return nil, fmt.Errorf("threshold exceeds number of preferences for voter %d", ind)
		}
		for i := 0; i < thresholds[ind]; i++ {
			count[prefs[i]]++
		}
	}
	return count, nil
}

func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error) {
	count, err := ApprovalSWF(p, thresholds)
	if err != nil {
		return nil, err
	}
	return MaxCount(count), nil
}
