package tp3

import "tp3/comsoc"

type Request struct {
	Operator string `json:"op"`
	Args     [2]int `json:"args"`
}

type Response struct {
	Result int `json:"res"`
}

type NewBallotRequest struct {
	Rule     string               `json:"rule"`
	Deadline string               `json:"deadline"`
	VoterIds []string             `json:"voter-ids"`
	Alts     int64                `json:"#alts"`
	TieBreak []comsoc.Alternative `json:"tie-break"`
}

type VoteRequest struct {
	AgentID  string `json:"agent-id"`
	BallotID string `json:"ballot-id"`
	Prefs    []int  `json:"prefs"`
	//options  []comsoc.Alternative `json:"tie-break"`
}

type NewBallotResponse struct {
	BallotID string `json:"ballot-id"`
}
