package agt

import (
	"sync"
	"time"
	"tp3/comsoc"
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
