package tp3

import "tp3/comsoc"

type NewBallotRequest struct {
	Rule     string               `json:"rule"`
	Deadline string               `json:"deadline"`
	VoterIds []string             `json:"voter-ids"`
	Alts     int64                `json:"#alts"`
	TieBreak []comsoc.Alternative `json:"tie-break"`
}

type VoteRequest struct {
	AgentID  string               `json:"agent-id"`
	BallotID string               `json:"ballot-id"`
	Prefs    []comsoc.Alternative `json:"prefs"`
	Options  []int                `json:"options"`
}

type ResultsRequest struct {
	BallotID string `json:"ballot-id"`
	//options  []comsoc.Alternative `json:"tie-break"`
}

type NewBallotResponse struct {
	BallotID string `json:"ballot-id"`
}

type ResultResponse struct {
	Winner  comsoc.Alternative   `json:"winner"`
	Ranking []comsoc.Alternative `json:"ranking"`
}
