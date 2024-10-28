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

func NewBallotAgent(ballotID string, rulename string, deadline time.Time, voterID map[string]bool, alts []comsoc.Alternative, tiebreak []comsoc.Alternative) *ballotAgent {
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

func (*ServerRestAgent) decodeRequestResult(r *http.Request) (req rad.ResultsRequest, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *ServerRestAgent) newBallotRest(w http.ResponseWriter, r *http.Request) {

	// vérification de la méthode de la requête ->ici on veut un POST
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestBallot(r)
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
	voterIDMap := make(map[string]bool)
	for _, name := range req.VoterIds {
		voterIDMap[name] = false
	}

	newBallot := *NewBallotAgent(ballotName, req.Rule, deadline, voterIDMap, make([]comsoc.Alternative, 0), req.TieBreak)
	for i := int64(0); i < req.Alts; i++ {
		newBallot.alternatives = append(newBallot.alternatives, comsoc.Alternative(i))
	}
	err = comsoc.CheckProfile(req.TieBreak, newBallot.alternatives)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Erreur dans la formation du tie-break")
		return
	}
	log.Println("laaaa")

	switch req.Rule {
	case "majority":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.MajoritySWF, comsoc.TieBreakFactory(req.TieBreak))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.MajoritySCF, comsoc.TieBreakFactory(req.TieBreak))
	case "borda":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.BordaSWF, comsoc.TieBreakFactory(req.TieBreak))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.BordaSCF, comsoc.TieBreakFactory(req.TieBreak))
	case "copeland":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.CopelandSWF, comsoc.TieBreakFactory(req.TieBreak))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.CopelandSCF, comsoc.TieBreakFactory(req.TieBreak))
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Méthode de vote non connu")
		return
	}
	rsa.ballotAgents[newBallot.ballotID] = &newBallot

	var resp rad.NewBallotResponse
	resp.BallotID = ballotName
	rsa.count++
	w.WriteHeader(http.StatusCreated)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}

func (rsa *ServerRestAgent) vote(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête ->ici on veut un POST
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestVote(r)
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
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	var ballotWanted *ballotAgent = rsa.ballotAgents[req.BallotID]

	if ballotWanted.deadline.Before(time.Now()) {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	//Regarder si l'agent peut voter
	log.Println("agent : ", req.AgentID)
	if value, exists := ballotWanted.voterID[req.AgentID]; exists {
		if value == true {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	} else {
		fmt.Fprint(w, "L'agent n'est pas autorisé à voter")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Regarder si l'agent a bien donné ses préférences
	if comsoc.CheckProfile(req.Prefs, ballotWanted.alternatives) != nil {
		fmt.Fprint(w, "Le bulletin n'est pas correcte")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Ajour les preferences du votant au Profile du vote
	ballotWanted.profile = append(ballotWanted.profile, req.Prefs)
	//Indiquer que l'agent a voté
	ballotWanted.voterID[req.AgentID] = true
	w.WriteHeader(http.StatusOK)

}

func (rsa *ServerRestAgent) results(w http.ResponseWriter, r *http.Request) {
	// vérification de la méthode de la requête ->ici on veut un POST
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := rsa.decodeRequestResult(r)
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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var ballotWanted *ballotAgent = rsa.ballotAgents[req.BallotID]

	if len(ballotWanted.profile) == 0 {
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprint(w, "Personne n'a voté :(")

		return
	}

	if ballotWanted.deadline.After(time.Now()) {
		w.WriteHeader(http.StatusTooEarly)
		return
	}

	ranking, _ := ballotWanted.ruleSWF(ballotWanted.profile)
	winner, _ := ballotWanted.ruleSCF(ballotWanted.profile)

	var resp rad.ResultResponse
	resp.Ranking = ranking
	resp.Winner = winner

	serial, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(serial)

}

func (rsa *ServerRestAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/new_ballot", rsa.newBallotRest)
	mux.HandleFunc("/vote", rsa.vote)
	mux.HandleFunc("/results", rsa.results)

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