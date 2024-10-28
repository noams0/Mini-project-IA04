package main

import (
	"log"
	"time"
	agt "tp3/agt"
	"tp3/comsoc"
)

func main() {
	administrator := agt.NewAdmin("adminAgent")
	deadline := time.Now().Add(4 * time.Second).Format(time.RFC3339)
	alts := make([]comsoc.Alternative, 6)
	for i := 1; i <= 6; i++ {
		alts[i-1] = comsoc.Alternative(i)
	}
	voterIDs := []string{"agt1", "agt2", "agt3"}
	ballotID, err := administrator.StartSession("majority", deadline, voterIDs, 6, alts)
	log.Println(err, ballotID)
}
