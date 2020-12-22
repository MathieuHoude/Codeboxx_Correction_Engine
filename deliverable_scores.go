package main

import (
	"encoding/json"
	"strings"
)

//DeliverableScore represents a deliverable criteria
type DeliverableScore struct {
	ID                int
	ScoreCardItemName string
	Pass              bool
}

func buildDeliverableScoresFromRspec(testResults RspecResults, deliverablesScores []DeliverableScore) []DeliverableScore {

	for _, testResult := range testResults.Examples {
		for i := 0; i < len(deliverablesScores); i++ {
			if testResult.Description == deliverablesScores[i].ScoreCardItemName {
				if testResult.Status == "passed" {
					deliverablesScores[i].Pass = true
				}
				break
			}
		}

	}
	return deliverablesScores
}

func buildDeliverableScores(content []byte, gradingRequest GradingRequest) []DeliverableScore {
	deliverableScores := gradingRequest.DeliverableScores

	switch gradingRequest.TestingTool {
	case "Rspec":
		var rspecResults RspecResults
		if err := json.NewDecoder(strings.NewReader(string(content))).Decode(&rspecResults); err != nil {
			panic(err)
		}
		deliverableScores = buildDeliverableScoresFromRspec(rspecResults, deliverableScores)
	}
	return deliverableScores
}
