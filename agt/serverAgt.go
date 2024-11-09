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

func NewBallotAgent(ballotID string, rulename string, deadline time.Time, voterID map[string]bool, alts []comsoc.Alternative, tiebreak []comsoc.Alternative, thresholds []int) *ballotAgent {

	if thresholds == nil {
		thresholds = []int{} // Crée une slice vide si non fournis
	}

	return &ballotAgent{ballotID: ballotID, rulename: rulename, deadline: deadline, voterID: voterID, alternatives: alts, tiebreak: tiebreak, thresholds: thresholds}
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

	// Création de la deadline
	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Deadline dans le passé
	if time.Now().After(deadline) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rsa.Lock()
	defer rsa.Unlock()
	ballotName := fmt.Sprintf("scurtinNum%d", rsa.count)
	voterIDMap := make(map[string]bool)

	// Si met 2 fois le même agent ne pose pas de problème
	for _, name := range req.VoterIds {
		voterIDMap[name] = false
	}

	// Transformer le tie-break de int en alternatives
	TieBreakAlt := comsoc.IntSliceToAlternativeSlice(req.TieBreak)

	// Créer le newBallot mais on ne connait pas les tresholds s'il y en a, ils seront passé pendant le vote
	// Pour l'instant on initialise à un slice vide
	newBallot := *NewBallotAgent(ballotName, req.Rule, deadline, voterIDMap, make([]comsoc.Alternative, 0), TieBreakAlt, nil)

	for i := int(0); i < req.Alts; i++ {
		newBallot.alternatives = append(newBallot.alternatives, comsoc.Alternative(i))
	}

	err = comsoc.CheckProfile(TieBreakAlt, newBallot.alternatives)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Erreur dans la formation du tie-break")
		return
	}

	switch req.Rule {
	case "majority":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.MajoritySWF, comsoc.TieBreakFactory(TieBreakAlt))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.MajoritySCF, comsoc.TieBreakFactory(TieBreakAlt))
	case "borda":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.BordaSWF, comsoc.TieBreakFactory(TieBreakAlt))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.BordaSCF, comsoc.TieBreakFactory(TieBreakAlt))
	case "copeland":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.CopelandSWF, comsoc.TieBreakFactory(TieBreakAlt))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.CopelandSCF, comsoc.TieBreakFactory(TieBreakAlt))
	case "approval":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.ApprovalSWF, comsoc.TieBreakFactory(TieBreakAlt))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.ApprovalSCF, comsoc.TieBreakFactory(TieBreakAlt))
	case "condorcet":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.CopelandSWF, comsoc.TieBreakFactory(TieBreakAlt))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.CondorcetWinner, comsoc.TieBreakFactory(TieBreakAlt))
	case "stv":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.StvSWF, comsoc.TieBreakFactory(TieBreakAlt))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.StvSCF, comsoc.TieBreakFactory(TieBreakAlt))
	case "kemeny":
		newBallot.ruleSWF = comsoc.SWFFactory(comsoc.KemenySWF, comsoc.TieBreakFactory(TieBreakAlt))
		newBallot.ruleSCF = comsoc.SCFFactory(comsoc.KemenySCF, comsoc.TieBreakFactory(TieBreakAlt))

	default:
		w.WriteHeader(http.StatusNotImplemented)
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
		if value {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "L'agent n'est pas autorisé à voter")
		return
	}

	// Transforme les préférences en alternatives
	PrefsAlt := comsoc.IntSliceToAlternativeSlice(req.Prefs)

	//Regarder si l'agent a bien donné ses préférences
	if comsoc.CheckProfile(PrefsAlt, ballotWanted.alternatives) != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Le bulletin n'est pas correcte")
		return
	}

	//Ajour les preferences du votant au Profile du vote
	ballotWanted.profile = append(ballotWanted.profile, PrefsAlt)

	// Ajouts des options uniquement pour le approval.
	var treshold int
	if req.Options != nil && len(req.Options) != 0 {

		// Vérification de req.Options de 0
		// doit être positif et ne pas être supérieur aux nombres d'alternatives
		if req.Options[0] > 0 && req.Options[0] <= len(ballotWanted.alternatives) {
			treshold = int(req.Options[0])
		}
	} // Si pas d'options ne vote pas

	ballotWanted.thresholds = append(ballotWanted.thresholds, treshold)

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

	var ranking []comsoc.Alternative
	var winner comsoc.Alternative
	if ballotWanted.rulename != "approval" {
		ranking, _ = ballotWanted.ruleSWF(ballotWanted.profile, comsoc.TransformInt(ballotWanted.tiebreak))
		winner, _ = ballotWanted.ruleSCF(ballotWanted.profile, comsoc.TransformInt(ballotWanted.tiebreak))
	} else { // Si approval met les thresholds
		ranking, _ = ballotWanted.ruleSWF(ballotWanted.profile, ballotWanted.thresholds)
		winner, _ = ballotWanted.ruleSCF(ballotWanted.profile, ballotWanted.thresholds)
	}

	var resp rad.ResultResponse
	resp.Ranking = comsoc.TransformInt(ranking)
	resp.Winner = int(winner)

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
