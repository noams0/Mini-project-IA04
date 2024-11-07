package main

import (
	"fmt"
	"time"
	agt "tp3/agt"
	"tp3/comsoc"
)

func main() {

	/*creating new ballot*/
	administrator := agt.NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)

	//Tie-break du plus petit au plus gros
	tb := make([]comsoc.Alternative, 6)
	for i := 1; i <= 6; i++ {
		tb[i-1] = comsoc.Alternative(i)
	}
	fmt.Println(tb)

	voterIDs := []string{"agt1", "agt2", "agt3", "agt4", "agt5", "agt6"}
	ballotIDs := []string{}
	list_voter := []*agt.Agent{}
	/*creating voting Agent*/
	agt1_4_preferences := []comsoc.Alternative{1, 5, 2, 4, 3, 6}
	agt2_preferences := []comsoc.Alternative{6, 4, 3, 5, 2, 1}
	agt3_preferences := []comsoc.Alternative{4, 3, 2, 6, 5, 1}
	agt5_preferences := []comsoc.Alternative{6, 3, 2, 1, 5, 4}
	agt1_4_options := []int{2}
	agt2_options := []int{1}
	agt3_options := []int{1}
	agt5_options := []int{1}

	agt1 := agt.NewAgent("agt1", agt1_4_preferences, agt1_4_options)
	list_voter = append(list_voter, agt1)
	agt2 := agt.NewAgent("agt2", agt2_preferences, agt2_options)
	list_voter = append(list_voter, agt2)
	agt3 := agt.NewAgent("agt3", agt3_preferences, agt3_options)
	list_voter = append(list_voter, agt3)
	agt4 := agt.NewAgent("agt4", agt1_4_preferences, agt1_4_options)
	list_voter = append(list_voter, agt4)
	agt5 := agt.NewAgent("agt5", agt5_preferences, agt5_options)
	list_voter = append(list_voter, agt5)

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
	//Point/Ranking -> 2, 1, 4, 6, 5 //OK
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

	ballotID, _ = administrator.StartSession("stv", deadline, voterIDs, 6, tb)
	// Point/Ranking -> 6,1,4,2,3,5 // Ok ranking mais pas resultat
	ballotIDs = append(ballotIDs, ballotID)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	time.Sleep(5 * time.Second)
	for _, ballot := range ballotIDs {
		administrator.GetResults(ballot)
	}

}
