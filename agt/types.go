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
	rule         func(comsoc.Profile) ([]comsoc.Alternative, error)
	deadline     time.Time
	voterID      map[string]bool
	profile      comsoc.Profile
	alternatives []comsoc.Alternative
	tiebreak     []comsoc.Alternative
	//thresholds   []int64
}
