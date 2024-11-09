package agt

import (
	"fmt"
	"testing"
	"time"
	"tp3/comsoc"
)

func TestGeneralExceptionsOptionsVote(t *testing.T) {
	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := make([]int, 6)
	for i := 1; i <= 6; i++ {
		tb[i-1] = i
	}

	voterIDs := []string{"agt1", "agt2", "agt3", "agt4", "agt5", "agt6", "agt7"}
	ballotIDs := []string{}
	list_voter := []*Agent{}
	/*creating voting Agent*/

	// vote en plus, pas bon nombre
	agt1_4_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6, 7}
	// Affiche bien bad request !

	// vote pas de 1 à 6
	agt2_preferences := []comsoc.Alternative{6, 4, 3, 5, 2, 7}
	//  Affiche bien bad request !

	// Vote vide
	// agt3_preferences := []comsoc.Alternative{}
	// Affiche bien bad request

	// deux fois la même alternative
	agt3_preferences := []comsoc.Alternative{6, 5, 2, 4, 3, 6}
	// Affiche bien bad request

	agt5_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt6_preferences := []comsoc.Alternative{1, 3, 2, 6, 5, 4}

	agt7_preferences := []comsoc.Alternative{1, 3, 2, 6, 5, 4}

	// mises d'options à nil
	agt1_4_options := []int{2}
	agt2_options := []int{1}
	agt3_options := []int{1}
	agt5_options := []int{1}

	// test avec options vide
	agt6_options := []int{}
	// marche empêche le vote de 6

	// test avec options au dessus du nombre d'alternatives
	// agt7_options := []int{9}
	// marche empêche le vote de 6

	// test avec plusieurs nombres dans les options
	agt7_options := []int{2, 4, 1}
	// Marche ne prend que la première des options

	agt1 := NewAgent("agt1", agt1_4_preferences, agt1_4_options)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, agt2_options)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, agt3_options)
	list_voter = append(list_voter, agt3)
	agt4 := NewAgent("agt4", agt1_4_preferences, agt1_4_options)
	list_voter = append(list_voter, agt4)
	agt5 := NewAgent("agt5", agt5_preferences, agt5_options)
	list_voter = append(list_voter, agt5)
	agt6 := NewAgent("agt6", agt6_preferences, agt6_options)
	list_voter = append(list_voter, agt6)

	// vote avec nil en options
	// Non pris en compte
	// agt7 := NewAgent("agt7", agt7_preferences, nil)
	agt7 := NewAgent("agt7", agt7_preferences, agt7_options)

	list_voter = append(list_voter, agt7)

	ballotID, _ := administrator.StartSession("majority", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK
	ballotIDs = append(ballotIDs, ballotID)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	//Point/Ranking ->2 : 19 / 3 : 17 / 4: 15 / 1 : 12 / 6 :  8/  5 : 6 OK
	ballotID, _ = administrator.StartSession("borda", deadline, voterIDs, 6, tb)
	ballotIDs = append(ballotIDs, ballotID)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	//Point/Ranking -> 2 : 5 / 3 : 3 / 1 : -1  / 4 : -1 / 5 : -3 / 6 : -5 OK
	ballotID, _ = administrator.StartSession("copeland", deadline, voterIDs, 6, tb)
	ballotIDs = append(ballotIDs, ballotID)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	ballotID, _ = administrator.StartSession("approval", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1,5,6,4,2,3 // OK
	ballotIDs = append(ballotIDs, ballotID)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	ballotID, _ = administrator.StartSession("condorcet", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 2, 1, 4, 6, 5 //OK
	ballotIDs = append(ballotIDs, ballotID)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}
	ballotID, _ = administrator.StartSession("stv", deadline, voterIDs, 6, tb)
	ballotIDs = append(ballotIDs, ballotID)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	ballotID, _ = administrator.StartSession("kemeny", deadline, voterIDs, 6, tb)
	// Point/Ranking -> 6,1,4,2,3,5
	ballotIDs = append(ballotIDs, ballotID)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	time.Sleep(5 * time.Second)
	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}
}

func TestWrongTieBreak(t *testing.T) {

	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := []int{2, 3, 4, 5, 6, 7}

	voterIDs := []string{"agt1", "agt2", "agt3"}
	ballotIDs := []string{}
	list_voter := []*Agent{}

	/*creating voting Agent*/

	agt1_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}

	agt2_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt3_preferences := []comsoc.Alternative{1, 3, 2, 6, 5, 4}

	agt1 := NewAgent("agt1", agt1_preferences, nil)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, nil)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, nil)
	list_voter = append(list_voter, agt3)

	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK

	ballotIDs = append(ballotIDs, ballotID)

	if err != nil {
		fmt.Println(err.Error()) // affiche bien l'erreur de tie break
		return
	}

	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	time.Sleep(5 * time.Second)
	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

}

func TestWrongDeadline(t *testing.T) {

	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(-4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := []int{1, 2, 3, 4, 5, 6}

	voterIDs := []string{"agt1", "agt2", "agt3"}
	ballotIDs := []string{}
	list_voter := []*Agent{}

	/*creating voting Agent*/

	agt1_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}

	agt2_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt3_preferences := []comsoc.Alternative{1, 3, 2, 6, 5, 4}

	agt1 := NewAgent("agt1", agt1_preferences, nil)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, nil)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, nil)
	list_voter = append(list_voter, agt3)

	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK

	ballotIDs = append(ballotIDs, ballotID)

	if err != nil {
		fmt.Println(err.Error()) // affiche bien l'erreur de deadline trop tôt
		return
	}

	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	time.Sleep(5 * time.Second)
	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

}

func TestVoteToLate(t *testing.T) {

	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := []int{1, 2, 3, 4, 5, 6}

	voterIDs := []string{"agt1", "agt2", "agt3"}
	ballotIDs := []string{}
	list_voter := []*Agent{}

	/*creating voting Agent*/

	agt1_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}

	agt2_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt3_preferences := []comsoc.Alternative{1, 3, 2, 6, 5, 4}

	agt1 := NewAgent("agt1", agt1_preferences, nil)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, nil)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, nil)
	list_voter = append(list_voter, agt3)

	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK

	ballotIDs = append(ballotIDs, ballotID)

	if err != nil {
		fmt.Println(err.Error()) // affiche bien l'erreur de deadline trop tôt
		return
	}

	for i, ag := range list_voter {
		if i != len(list_voter)-1 {
			ag.Vote(ballotID)
		}
	}

	time.Sleep(5 * time.Second)

	// Vote trop tard
	list_voter[len(list_voter)-1].Vote(ballotID)
	// Bien pas pris en compte !

	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

}

func TestVoteNotAuthorized(t *testing.T) {

	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := []int{1, 2, 3, 4, 5, 6}

	voterIDs := []string{"agt1", "agt2"}
	ballotIDs := []string{}
	list_voter := []*Agent{}

	/*creating voting Agent*/

	agt1_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}

	agt2_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt3_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt1 := NewAgent("agt1", agt1_preferences, nil)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, nil)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, nil)
	list_voter = append(list_voter, agt3)

	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK

	ballotIDs = append(ballotIDs, ballotID)

	if err != nil {
		fmt.Println(err.Error()) // affiche bien l'erreur de deadline trop tôt
		return
	}

	// le vote de agt3 est bien pas pris en compte
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	time.Sleep(5 * time.Second)

	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

}

func TestResultsToEarly(t *testing.T) {

	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := []int{1, 2, 3, 4, 5, 6}

	voterIDs := []string{"agt1", "agt2", "agt3"}
	ballotIDs := []string{}
	list_voter := []*Agent{}

	/*creating voting Agent*/

	agt1_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}

	agt2_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt3_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt1 := NewAgent("agt1", agt1_preferences, nil)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, nil)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, nil)
	list_voter = append(list_voter, agt3)

	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK

	ballotIDs = append(ballotIDs, ballotID)

	if err != nil {
		fmt.Println(err.Error()) // affiche bien l'erreur de deadline trop tôt
		return
	}

	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	// Attend pas assez avant le résultat
	time.Sleep(2 * time.Second)

	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}
	// Affiche bien l'erreur

	// Attends assez
	time.Sleep(2 * time.Second)

	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

}

func TestWrongRule(t *testing.T) {

	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := []int{1, 2, 3, 4, 5, 6}

	voterIDs := []string{"agt1", "agt2", "agt3"}
	ballotIDs := []string{}
	list_voter := []*Agent{}

	/*creating voting Agent*/

	agt1_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}

	agt2_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt3_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt1 := NewAgent("agt1", agt1_preferences, nil)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, nil)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, nil)
	list_voter = append(list_voter, agt3)

	ballotID, err := administrator.StartSession("jhonny", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK
	if err != nil {
		fmt.Println(err) // affiche bien l'erreur de deadline trop tôt
		return
	}

	ballotIDs = append(ballotIDs, ballotID)

	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	// Attend pas assez avant le résultat
	time.Sleep(2 * time.Second)

	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}
	// Affiche bien l'erreur

	// Attends assez
	time.Sleep(2 * time.Second)

	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

}

func TestVoteDejaEffectuerEtMauvaisBallotEtTropTard(t *testing.T) {

	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := []int{1, 2, 3, 4, 5, 6}

	voterIDs := []string{"agt1", "agt2", "agt3"}
	ballotIDs := []string{}
	list_voter := []*Agent{}

	/*creating voting Agent*/

	agt1_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}

	agt2_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt3_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt1 := NewAgent("agt1", agt1_preferences, nil)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, nil)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, nil)
	list_voter = append(list_voter, agt3)

	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK
	if err != nil {
		fmt.Println(err) // affiche bien l'erreur de deadline trop tôt
		return
	}

	ballotIDs = append(ballotIDs, ballotID)

	for i, ag := range list_voter {
		if i != len(list_voter)-1 {
			ag.Vote(ballotID)
		}

	}

	// Revote
	list_voter[0].Vote(ballotID)

	// vote dans ballot pas implémentés
	//list_voter[len(list_voter)-1].Vote("PRairIe")

	time.Sleep(5 * time.Second)

	// Vote trop tard pour agent 3
	list_voter[len(list_voter)-1].Vote(ballotID)

	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

}

func TestResultNotFound(t *testing.T) {

	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := []int{1, 2, 3, 4, 5, 6}

	voterIDs := []string{"agt1", "agt2", "agt3"}
	ballotIDs := []string{}
	list_voter := []*Agent{}

	/*creating voting Agent*/

	agt1_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}

	agt2_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt3_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt1 := NewAgent("agt1", agt1_preferences, nil)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, nil)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, nil)
	list_voter = append(list_voter, agt3)

	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK
	if err != nil {
		fmt.Println(err) // affiche bien l'erreur de deadline trop tôt
		return
	}

	ballotIDs = append(ballotIDs, ballotID)

	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	time.Sleep(5 * time.Second)

	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

	administrator.GetResults("BallotInexistant")

}

func TestMemeVotantsMultiples(t *testing.T) {

	// Lauching the server
	server := NewServerRestAgent(":8080")
	go server.Start()

	/*creating new ballot*/
	administrator := NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := []int{1, 2, 3, 4, 5, 6}

	// erreur dans la liste des agents 2 fois le même
	// Ne sera pas capable de revoter !
	voterIDs := []string{"agt1", "agt2", "agt3", "agt1"}
	ballotIDs := []string{}
	list_voter := []*Agent{}

	/*creating voting Agent*/

	agt1_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}

	agt2_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt3_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}

	agt1 := NewAgent("agt1", agt1_preferences, nil)
	list_voter = append(list_voter, agt1)
	agt2 := NewAgent("agt2", agt2_preferences, nil)
	list_voter = append(list_voter, agt2)
	agt3 := NewAgent("agt3", agt3_preferences, nil)
	list_voter = append(list_voter, agt3)

	// Deux fois le même agent dans les votes
	// Bien pas repris en compte !
	list_voter = append(list_voter, agt1)

	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, tb)
	//Point/Ranking -> 1, 2, 4, 6, 3, 5 OK
	if err != nil {
		fmt.Println(err) // affiche bien l'erreur de deadline trop tôt
		return
	}

	ballotIDs = append(ballotIDs, ballotID)

	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	time.Sleep(5 * time.Second)

	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

}
