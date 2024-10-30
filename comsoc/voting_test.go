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

	res, err := CondorcetWinner(prefs)

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
