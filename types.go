package tp3

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
