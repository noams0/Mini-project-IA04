package agt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	rad "tp3"
	"tp3/comsoc"
)

func NewAgent(id string, prefs []comsoc.Alternative, options []int) *Agent {

	if options == nil {
		options = []int{}
	}

	return &Agent{id, prefs, options}
}

func NewAdmin(id string) *Admin {
	return &Admin{id}
}

func (ag Agent) Clone() *Agent {
	return NewAgent(ag.agentId, ag.prefs, ag.options)
}

func (ag Agent) String() string {
	return fmt.Sprintf("ID : %s, Preferences : %v, Options : %v", ag.agentId, ag.prefs, ag.options)
}

func (ad Admin) DecodeNewBallotResponse(r *http.Response) (rad.NewBallotResponse, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return rad.NewBallotResponse{}, err
	}

	var resp rad.NewBallotResponse

	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		fmt.Println("failed unmarshalling")
		return rad.NewBallotResponse{}, err
	}

	return resp, nil
}

func (ad Admin) StartSession(rule string, deadline string, voterIds []string, alts int64, tieBreak []comsoc.Alternative) (res string, err error) {

	requestURL := "http://localhost:8080/new_ballot"

	session := rad.NewBallotRequest{
		Rule:     rule,
		Deadline: deadline,
		VoterIds: voterIds,
		Alts:     alts,
		TieBreak: tieBreak,
	}

	log.Println(session)
	data, _ := json.Marshal(session)
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(data))

	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		fmt.Println(err)
		return
	}

	result, err := ad.DecodeNewBallotResponse(resp)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed treating response")
		return
	}

	return result.BallotID, nil
}

func (ag Agent) Vote(sessionID string) {
	requestURL := "http://localhost:8080/vote"

	vote := rad.VoteRequest{
		AgentID:  ag.agentId,
		BallotID: sessionID,
		Prefs:    ag.prefs,
		Options:  ag.options,
	}

	data, _ := json.Marshal(vote)

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		fmt.Println(err)
		return
	}

	fmt.Println("Vote bien pris en compte")
	return
}

func (ad Admin) GetResults(sessionID string) {

	requestURL := "http://localhost:8080/results"
	obj := rad.ResultsRequest{BallotID: sessionID}
	data, _ := json.Marshal(obj)

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	_, err2 := buf.ReadFrom(resp.Body)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	var result rad.ResultResponse
	result.Ranking = make([]comsoc.Alternative, 0)
	err = json.Unmarshal(buf.Bytes(), &result)
	fmt.Println("resultat", result)

	if err != nil {
		fmt.Println("failed unmarshalling")
		return
	}
	fmt.Println()

	if result.Winner == -1 {
		fmt.Printf("Pas de gagnant de Condorcet pour le vote %s\n", sessionID)
	} else {
		fmt.Printf("Le gagnant du vote %s est %d\n ", sessionID, result.Winner)
	}
	if len(result.Ranking) > 0 {
		fmt.Printf("Et le classement est %v", result.Ranking)
		fmt.Println()
	}
}
