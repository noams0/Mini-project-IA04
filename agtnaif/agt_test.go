package agtnaif

import (
	"testing"
	"time"
	"tp3/comsoc"
)

// TestEqual vérifie la méthode Equal de la structure Agent
func TestEqual(t *testing.T) {
	agent1 := &Agent{
		ID:    1,
		Name:  "Agent 1",
		Prefs: []comsoc.Alternative{1, 2, 3},
	}

	agent2 := &Agent{
		ID:    2,
		Name:  "Agent 1",
		Prefs: []comsoc.Alternative{3, 2, 1},
	}

	agent3 := &Agent{
		ID:    3,
		Name:  "Agent 3",
		Prefs: []comsoc.Alternative{1, 3, 2},
	}

	if !agent1.Equal(agent2) {
		t.Errorf("Expected agent1 and agent2 to be equal, but they are not")
	}

	if agent1.Equal(agent3) {
		t.Errorf("Expected agent1 and agent3 to be different, but they are equal")
	}
}

// TestDeepEqual vérifie la méthode DeepEqual de la structure Agent
func TestDeepEqual(t *testing.T) {
	agent1 := &Agent{
		ID:    1,
		Name:  "Agent 1",
		Prefs: []comsoc.Alternative{1, 2, 3},
	}

	agent2 := &Agent{
		ID:    2,
		Name:  "Agent 1",
		Prefs: []comsoc.Alternative{1, 2, 3},
	}

	agent3 := &Agent{
		ID:    1,
		Name:  "Agent 1",
		Prefs: []comsoc.Alternative{1, 2, 3},
	}

	if !agent1.DeepEqual(agent3) {
		t.Errorf("Expected agent1 and agent3 to be deeply equal, but they are not")
	}

	if agent1.DeepEqual(agent2) {
		t.Errorf("Expected agent1 and agent2 to be deeply different, but they are equal")
	}
}

func TestAgentVoting(t *testing.T) {

	alternatives := []comsoc.Alternative{1, 2, 3}
	server := &VotingServer{
		C: make(chan VotingMessage),
		p: make(comsoc.Profile, 0),
		Scf: comsoc.SCFFactory(
			comsoc.MajoritySCF,
			comsoc.TieBreakFactory([]comsoc.Alternative{1, 2, 3, 4}),
		),
	}

	// Create agents
	agent1 := &Agent{
		ID:         1,
		Name:       "Agent 1",
		Prefs:      alternatives,
		ServerChan: server.C,
	}
	agent2 := &Agent{
		ID:         2,
		Name:       "Agent 2",
		Prefs:      []comsoc.Alternative{2, 3, 1},
		ServerChan: server.C,
	}

	// Start the server in a goroutine
	go server.StartServer()
	defer server.StopServer()

	// Start a vote on the server
	server.StartVoteTime(10 * time.Second) // Vote open for 5 seconds

	// Start the agents in their own goroutines
	go agent1.Start()
	go agent2.Start()
	//fmt.Println("here")

	// Give the agents time to vote and the server time to close the vote
	time.Sleep(10 * time.Second)

	// Check if the winner was determined correctly
	if server.lastWinner != 1 {
		t.Errorf("Expected winner to be 1, but got %d", server.lastWinner)
	}

	// Check if the profiles were updated correctly
	if len(server.p) != 2 {
		t.Errorf("Expected 2 profiles in server, but got %d", len(server.p))
	}

	if server.p[0][0] != 1 || server.p[1][0] != 2 {
		t.Errorf("Unexpected preferences in server profile")
	}
}

func TestAgentClone(t *testing.T) {
	// Test the Clone function
	agent := &Agent{
		ID:    1,
		Name:  "Agent 1",
		Prefs: []comsoc.Alternative{1, 2, 3},
	}

	clonedAgent := agent.Clone()

	if !agent.DeepEqual(clonedAgent) {
		t.Errorf("Expected agent and cloned agent to be equal, but they are not")
	}
}
