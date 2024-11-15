package agt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/noams0/Mini-project-IA04/comsoc"
	"log"
	"net/http"
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

func (ad Admin) DecodeNewBallotResponse(r *http.Response) (NewBallotResponse, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return NewBallotResponse{}, err
	}

	var resp NewBallotResponse

	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		fmt.Println("failed unmarshalling")
		return NewBallotResponse{}, err
	}

	return resp, nil
}

func (ad Admin) StartVotingSession(rule string, deadline string, voterIds []string, alts int, tieBreak []int) (res string, err error) {

	requestURL := "http://localhost:8080/new_ballot"

	session := NewBallotRequest{
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

	PrefsInt := comsoc.TransformInt(ag.prefs)

	vote := VoteRequest{
		AgentID:  ag.agentId,
		BallotID: sessionID,
		Prefs:    PrefsInt,
		Options:  ag.options,
	}

	data, _ := json.Marshal(vote)

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

	fmt.Println("Vote bien pris en compte")
	return
}

func (ad Admin) GetResults(sessionID string) {

	requestURL := "http://localhost:8080/results"
	obj := ResultsRequest{BallotID: sessionID}
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

	var result ResultResponse
	result.Ranking = make([]int, 0)
	err = json.Unmarshal(buf.Bytes(), &result)

	if err != nil {
		fmt.Println("failed unmarshalling")
		return
	}

	if result.Winner == -1 {
		fmt.Printf("\nPas de gagnant de Condorcet pour le vote %s\n", sessionID)
	} else {
		fmt.Printf("\nLe gagnant du vote %s est %d\n ", sessionID, result.Winner)
	}
	if len(result.Ranking) > 0 {
		fmt.Printf("Et le classement est %v\n", result.Ranking)
		fmt.Println()
	}
}
