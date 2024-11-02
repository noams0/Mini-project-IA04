package main

import (
	"log"
	"time"
	agt "tp3/agt"
	"tp3/comsoc"
)

func main() {

	/*creating new ballot*/
	administrator := agt.NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)
	alts := make([]comsoc.Alternative, 6)
	for i := 1; i <= 6; i++ {
		alts[i-1] = comsoc.Alternative(i)
	}

	voterIDs := []string{"agt1", "agt2", "agt3", "agt4", "agt5", "agt6"}
	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, alts)
	log.Println(ballotID, err)
	list_voter := []*agt.Agent{}
	/*creating voting Agent*/
	agt1_4_preferences := []comsoc.Alternative{1, 2, 3, 4, 5, 6}
	agt2_preferences := []comsoc.Alternative{2, 4, 3, 5, 6, 1}
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

	for _, ag := range list_voter {
		ag.Vote("scurtinNum0")
	}

	ballotID, _ = administrator.StartSession("borda", deadline, voterIDs, 6, alts)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}
	ballotID, _ = administrator.StartSession("copeland", deadline, voterIDs, 6, alts)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	ballotID, _ = administrator.StartSession("approval", deadline, voterIDs, 6, alts)
	for _, ag := range list_voter {
		ag.Vote(ballotID)
	}

	time.Sleep(5 * time.Second)

	administrator.GetResults("scurtinNum0")
	administrator.GetResults("scurtinNum1")
	administrator.GetResults("scurtinNum2")
	administrator.GetResults("scurtinNum3")

}
