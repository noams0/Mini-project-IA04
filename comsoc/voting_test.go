package comsoc

import (
	"fmt"
	"testing"
)

func TestBordaSWF(t *testing.T) {

	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}

	res, _ := BordaSWF(prefs, nil)
	if res[1] != 4 {
		t.Errorf("error, result for 1 should be 4, %d computed", res[1])
	}
	if res[2] != 3 {
		t.Errorf("error, result for 2 should be 3, %d computed", res[2])
	}
	if res[3] != 2 {
		t.Errorf("error, result for 3 should be 2, %d computed", res[3])
	}
}

func TestBordaSCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	res, err := BordaSCF(prefs, nil)
	if err != nil {
		t.Error(err)
	}
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestMajoritySWF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	res, _ := MajoritySWF(prefs, nil)
	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 0 {
		t.Errorf("error, result for 2 should be 0, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestMajoritySCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	res, err := MajoritySCF(prefs, nil)
	if err != nil {
		t.Error(err)
	}
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestApprovalSWF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	thresholds := []int{2, 1, 2}
	res, _ := ApprovalSWF(prefs, thresholds)
	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[2] != 2 {
		t.Errorf("error, result for 2 should be 2, %d computed", res[2])
	}
	if res[3] != 1 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestApprovalSCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 3, 2},
		{1, 2, 3},
		{2, 1, 3},
	}
	thresholds := []int{2, 1, 2}
	res, err := ApprovalSCF(prefs, thresholds)

	if err != nil {
		t.Error(err)
	}
	if len(res) != 1 || res[0] != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestApprovalSCFWithTieBreak(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	thresholds := []int{2, 1, 2}
	bestAlts, err := ApprovalSCF(prefs, thresholds)
	if err != nil {
		t.Error(err)
	}
	if len(bestAlts) != 2 || bestAlts[0] != 1 || bestAlts[1] != 2 {
		t.Errorf("error, 1 and 2 should be the only best Alternatives")
	}
	orderedAlts := []Alternative{1, 2, 3}
	tieBreakFunc := TieBreakFactory(orderedAlts)
	winner, err := tieBreakFunc(bestAlts)
	if err != nil {
		t.Error(err)
	}
	if winner != 1 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestSWFFactory(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{2, 3, 1},
	}
	order := []Alternative{
		2, 3, 1,
	}
	tb := TieBreakFactory(order)

	newSwf := SWFFactory(BordaSWF, tb)
	res, _ := newSwf(prefs, nil)
	if res[0] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[0])
	}
}

func TestSWFFactoryApproval(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	thresholds := []int{2, 1, 2}

	order := []Alternative{
		2, 3, 1,
	}
	tb := TieBreakFactory(order)

	newSwf := SWFFactory(ApprovalSWF, tb)
	res, _ := newSwf(prefs, thresholds)

	fmt.Println(res)

	if res[0] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[1])
	}
	if res[1] != 1 {
		t.Errorf("error, result for 2 should be 2, %d computed", res[2])
	}
	if res[2] != 3 {
		t.Errorf("error, result for 3 should be 1, %d computed", res[3])
	}
}

func TestSCFFactory(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
		{3, 2, 1},
	}

	order := []Alternative{
		2, 3, 1,
	}
	tb := TieBreakFactory(order)

	newScf := SCFFactory(MajoritySCF, tb)

	alts, _ := newScf(prefs, nil)

	if alts != 3 {
		t.Errorf("error, result for 1 should be 3, %d computed", alts)
	}
}

func TestSCFFactoryApproval(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
	}
	thresholds := []int{2, 1, 2}
	order := []Alternative{
		2, 3, 1,
	}

	tb := TieBreakFactory(order)

	newScf := SCFFactory(ApprovalSCF, tb)

	res, err := newScf(prefs, thresholds)

	if err != nil {
		t.Error(err)
	}

	if res != 2 {
		t.Errorf("error, 1 should be the only best Alternative")
	}
}

func TestCondorcetWinner(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{3, 1, 2},
		{2, 1, 3},
	}

	res, err := CondorcetWinner(prefs, nil)

	if err != nil {
		t.Error(err)
	}

	if len(res) == 0 {
		t.Error("Pas de vainqueur Condorcet, alors que 1 est vainqueur de Condorcet")
		return
	}
	if len(res) != 1 {
		t.Error("Il ne devrait y avoir qu'un seul gagnant Condorcet,")
	}
	if res[0] != 1 {
		t.Errorf("error, result for 1 should be 1, %d computed", res[0])
	}
}

func TestCopelandSWF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{3, 1, 2},
		{2, 1, 3},
	}

	res, err := CopelandSWF(prefs, nil)
	if err != nil {
		t.Error(err)
	}

	if res[1] != 2 {
		t.Errorf("error, result for 1 should be 2, %d computed", res[0])
	}
	if res[2] != 1 {
		t.Errorf("error, result for 2 should be 1, %d computed", res[1])
	}
}

func TestCopelandSCF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{1, 3, 1},
		{3, 2, 1},
	}

	order := []Alternative{
		2, 3, 1,
	}
	tb := TieBreakFactory(order)

	newScf := SCFFactory(CopelandSCF, tb)

	alts, _ := newScf(prefs, nil)

	if alts != 1 {
		t.Errorf("error, result for 1 should be 1, %d computed", alts)
	}
}

func TestStvSwF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{1, 2, 3},
		{3, 2, 1},
		{3, 2, 1},
	}

	order := []Alternative{
		1, 3, 2,
	}

	order2 := []int{
		1, 3, 2,
	}

	tb := TieBreakFactory(order)

	newSwf := SWFFactory(StvSWF, tb)

	alts, _ := newSwf(prefs, order2)

	fmt.Println(alts)

}

func TestStvScF(t *testing.T) {
	prefs := [][]Alternative{
		{1, 5, 2, 4, 3, 6},
		{1, 5, 2, 4, 3, 6},
		{6, 4, 3, 5, 2, 1},
		{4, 3, 2, 6, 5, 1},
		{6, 3, 2, 1, 5, 4},
	}

	tb := make([]Alternative, 6)
	for i := 1; i <= 6; i++ {
		tb[i-1] = Alternative(i)
	}

	tbf := TieBreakFactory(tb)

	newScf := SCFFactory(StvSCF, tbf)

	alts, _ := newScf(prefs, TransformInt(tb))

	fmt.Println("resultat", alts)
}

func TestDistanceEdition(t *testing.T) {
	profile1 := []Alternative{1, 2, 3, 4}
	profile2 := []Alternative{1, 3, 4, 2}
	dist := DistanceEdition(profile1, profile2)
	if dist != 2 {
		t.Errorf("error, distance should be 2, %d computed", dist)
	}
}

func TestDistanceEditionProfile(t *testing.T) {
	prefs := [][]Alternative{
		{1, 2, 3},
		{3, 1, 2},
		{2, 1, 3},
	}
	rangement := []Alternative{1, 2, 3}
	dist := DistanceEditionProfile(prefs, rangement)
	if dist != 3 {
		t.Errorf("error, distance should be 3, %d computed", dist)
	}
	prefs = [][]Alternative{
		{1, 2, 3, 4},
		{4, 3, 1, 2},
		{2, 1, 4, 3},
	}
	rangement = []Alternative{1, 2, 3, 4}
	dist = DistanceEditionProfile(prefs, rangement)
	if dist != 7 {
		t.Errorf("error, distance should be 7, %d computed", dist)
	}
}

func TestKemenySWF(t *testing.T) {

	prefs := [][]Alternative{
		{1, 3, 2},
		{2, 1, 3},
	}
	//rangement_expected = []Alternative{2, 1, 3}
	rangement_expected := map[Alternative]int{
		1: 2,
		2: 3,
		3: 1,
	}
	rangement, _ := KemenySWF(prefs, []int{2, 3, 1})
	fmt.Println(rangement)
	for i, _ := range rangement {
		if rangement[i] != rangement_expected[i] {
			t.Errorf("error, for index %d expected %d, got %d", i, rangement_expected[i], rangement[i])
		}
	}

	// Nashville -> 1
	// Chattanooga -> 2
	// Knoxville -> 3
	// Memphis -> 4
	prefs = [][]Alternative{}
	//rangement_expected = []Alternative{1, 2, 3, 4}
	rangement_expected = map[Alternative]int{
		1: 4,
		2: 3,
		3: 2,
		4: 1,
	}
	for i := 0; i < 42; i++ {
		prefs = append(prefs, []Alternative{4, 1, 2, 3})
	}
	for i := 0; i < 26; i++ {
		prefs = append(prefs, []Alternative{1, 2, 3, 4})
	}
	for i := 0; i < 15; i++ {
		prefs = append(prefs, []Alternative{2, 3, 1, 4})
	}
	for i := 0; i < 17; i++ {
		prefs = append(prefs, []Alternative{3, 2, 1, 4})
	}
	rangement, _ = KemenySWF(prefs, nil)
	fmt.Println(rangement)
	for i, _ := range rangement {
		if rangement[i] != rangement_expected[i] {
			t.Errorf("error, for index %d expected %d, got %d", i, rangement_expected[i], rangement[i])
		}
	}
}

func TestKemenySCF(t *testing.T) {
	prefs := [][]Alternative{}

	for i := 0; i < 42; i++ {
		prefs = append(prefs, []Alternative{4, 1, 2, 3})
	}
	for i := 0; i < 26; i++ {
		prefs = append(prefs, []Alternative{1, 2, 3, 4})
	}
	for i := 0; i < 15; i++ {
		prefs = append(prefs, []Alternative{2, 3, 1, 4})
	}
	for i := 0; i < 17; i++ {
		prefs = append(prefs, []Alternative{3, 2, 1, 4})
	}
	order := []Alternative{
		1, 2, 3, 4,
	}
	tbf := TieBreakFactory(order)

	newScf := SCFFactory(KemenySCF, tbf)

	alts, _ := newScf(prefs, TransformInt(order))

	if alts != 1 {
		t.Errorf("error, result for 1 should be 1, %d computed", alts)
	}
}
