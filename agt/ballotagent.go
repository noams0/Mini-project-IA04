package agt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	rad "tp3"
	"tp3/comsoc"
)

func NewServerRestAgent(addr string) *ServerRestAgent {
	return &ServerRestAgent{id: addr, addr: addr, ballotAgents: make(map[string]*ballotAgent)}
}

func (rsa *ServerRestAgent) GetBallot() map[string]*ballotAgent {
	return rsa.ballotAgents
}

func NewBallotAgent(ballotID string, rulename string, deadline time.Time, voterID []string, alts []comsoc.Alternative, tiebreak []comsoc.Alternative) *ballotAgent {
	return &ballotAgent{ballotID: ballotID, rulename: rulename, deadline: deadline, voterID: voterID, alternatives: alts, tiebreak: tiebreak}
}

func (rsa *ServerRestAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

func (*ServerRestAgent) decodeRequestBallot(r *http.Request) (req rad.NewBallotRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (*ServerRestAgent) decodeRequestVote(r *http.Request) (req rad.VoteRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *ServerRestAgent) newBallotRest(w http.ResponseWriter, r *http.Request) {
	log.Println(rsa.GetBallot())
	// vérification de la méthode de la requête ->ici on veut un POST
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestBallot(r)
	log.Println(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if time.Now().After(deadline) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rsa.Lock()
	defer rsa.Unlock()
	ballotName := fmt.Sprintf("scurtinNum%d", rsa.count)
	log.Println(ballotName)
	newBallot := *NewBallotAgent(ballotName, req.Rule, deadline, req.VoterIds, make([]comsoc.Alternative, 0), req.TieBreak)
	for i := int64(0); i < req.Alts; i++ {
		newBallot.alternatives = append(newBallot.alternatives, comsoc.Alternative(i))
	}
	err = comsoc.CheckProfile(req.TieBreak, newBallot.alternatives)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch req.Rule {
	case "majority":
		newBallot.rule = comsoc.SWFFactory(comsoc.MajoritySWF, comsoc.TieBreakFactory(req.TieBreak))
	case "borda":
		newBallot.rule = comsoc.SWFFactory(comsoc.BordaSWF, comsoc.TieBreakFactory(req.TieBreak))
	case "copeland":
		newBallot.rule = comsoc.SWFFactory(comsoc.CopelandSWF, comsoc.TieBreakFactory(req.TieBreak))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rsa.ballotAgents[newBallot.ballotID] = &newBallot

	var resp rad.NewBallotResponse
	resp.BallotID = ballotName
	rsa.count++
	w.WriteHeader(http.StatusCreated)
	serial, _ := json.Marshal(resp) //la serialization qu'on a pas encore vu en cours
	w.Write(serial)
}

func (rsa *ServerRestAgent) vote(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête ->ici on veut un POST
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestVote(r)
	log.Println(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	balloPos := ""
	for i, b := range rsa.ballotAgents {
		if b.ballotID == req.BallotID {
			balloPos = i
		}
	}
	if balloPos == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (rsa *ServerRestAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.newBallotRest)
	mux.HandleFunc("/vote", rsa.vote)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
