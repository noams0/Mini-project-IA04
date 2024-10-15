package agtnaif

import (
	"fmt"
	"log"
	"time"
	"tp3/comsoc"
)

type VotingServerStateInformation struct {
	VoteOpen   bool
	LastWinner comsoc.Alternative
}

type VotingMessage struct {
	Preferences []comsoc.Alternative
	ReplyChan   chan VotingServerStateInformation
}

type VotingServer struct {
	C          chan VotingMessage
	p          comsoc.Profile
	voteOpen   bool
	serverOpen bool
	lastWinner comsoc.Alternative
	Scf        func(comsoc.Profile) (comsoc.Alternative, error)
}

func (vs *VotingServer) StartVoteTime(duration time.Duration) {
	vs.voteOpen = true
	log.Printf("vote open\n")
	go func(duration time.Duration) {
		time.Sleep(duration)
		vs.voteOpen = false
		log.Printf("vote closed\n")
		fmt.Println("error is comming : ", vs.p)
		vs.lastWinner, _ = vs.Scf(vs.p)
	}(duration)
}

func (vs *VotingServer) SetSCF(scf func(comsoc.Profile) (comsoc.Alternative, error)) {
	vs.Scf = scf
}

func (vs *VotingServer) StartServer() {
	vs.serverOpen = true
	for vs.serverOpen {
		select {
		case message := <-vs.C:
			if vs.voteOpen {
				vs.p = append(vs.p, message.Preferences)
			}
			message.ReplyChan <- VotingServerStateInformation{vs.voteOpen, vs.lastWinner}
		}
	}
}

func (vs *VotingServer) StopServer() {
	vs.serverOpen = false
}
