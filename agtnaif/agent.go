package agtnaif

import (
	"log"
	"reflect"
	"time"
	"tp3/comsoc"
)

type AgentI interface {
	Equal(ag AgentI) bool
	DeepEqual(ag AgentI) bool
	Clone() AgentI
	String() string
	Prefers(a comsoc.Alternative, b comsoc.Alternative) bool
	Start()
}

type Agent struct {
	ID         uint
	Name       string
	Prefs      []comsoc.Alternative
	ServerChan chan VotingMessage
}

func (agt *Agent) String() string {
	return agt.Name
}

func (agt *Agent) Equal(ag AgentI) bool {
	return agt.String() == ag.String()
}

func (agt *Agent) DeepEqual(ag AgentI) bool {
	if reflect.TypeOf(ag) != reflect.TypeOf(agt) {
		return false
	}

	return reflect.DeepEqual(agt, ag)
}

func (agt *Agent) Clone() AgentI {
	return &Agent{
		ID:         agt.ID,
		Name:       agt.Name,
		Prefs:      agt.Prefs,
		ServerChan: agt.ServerChan,
	}
}

func (agt *Agent) Prefers(a comsoc.Alternative, b comsoc.Alternative) bool {
	return comsoc.IsPref(a, b, agt.Prefs)
}

func (agt *Agent) Start() {
	replyChan := make(chan VotingServerStateInformation)

	// try to vote until you are successful -> not a dangerous method ?
	voted := false
	for !voted {
		time.Sleep(5 * time.Second)
		agt.ServerChan <- VotingMessage{
			Preferences: agt.Prefs,
			ReplyChan:   replyChan,
		}
		reply := <-replyChan
		if reply.VoteOpen {
			voted = true
		}
	}
	log.Printf("%s: voted\n", agt.Name)

	// wait for the vote to close to fetch the winner
	var winner comsoc.Alternative = -1
	for winner == -1 {
		time.Sleep(5 * time.Second)
		agt.ServerChan <- VotingMessage{
			Preferences: agt.Prefs,
			ReplyChan:   replyChan,
		}
		reply := <-replyChan
		if !reply.VoteOpen {
			winner = reply.LastWinner
		}
	}
	log.Printf("%s: %d won, it was in position %d in my preferences\n", agt.Name, winner, comsoc.Rank(winner, agt.Prefs))
}
