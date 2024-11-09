package agt

import (
	"github.com/noams0/Mini-project-IA04/comsoc"
	"sync"
	"time"
)

type ServerRestAgent struct {
	sync.Mutex
	id           string
	addr         string
	ballotAgents map[string]*ballotAgent
	count        int64
}

type ballotAgent struct {
	sync.Mutex
	ballotID     string
	rulename     string
	ruleSWF      func(comsoc.Profile, []int) ([]comsoc.Alternative, error)
	ruleSCF      func(comsoc.Profile, []int) (comsoc.Alternative, error)
	deadline     time.Time
	voterID      map[string]bool
	profile      comsoc.Profile
	alternatives []comsoc.Alternative
	tiebreak     []comsoc.Alternative
	thresholds   []int
}

type Agent struct {
	agentId string
	prefs   []comsoc.Alternative
	options []int
}

type Admin struct {
	agentId string
}

type NewBallotRequest struct {
	Rule     string   `json:"rule"`
	Deadline string   `json:"deadline"`
	VoterIds []string `json:"voter-ids"`
	Alts     int      `json:"#alts"`
	TieBreak []int    `json:"tie-break"`
}

type VoteRequest struct {
	AgentID  string `json:"agent-id"`
	BallotID string `json:"ballot-id"`
	Prefs    []int  `json:"prefs"`
	Options  []int  `json:"options"`
}

type ResultsRequest struct {
	BallotID string `json:"ballot-id"`
	//options  []comsoc.Alternative `json:"tie-break"`
}

type NewBallotResponse struct {
	BallotID string `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  int   `json:"winner"`
	Ranking []int `json:"ranking"`
}
