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

func buildDeliverableScoresFromJest(testResults JestResults, deliverablesScores []DeliverableScore) []DeliverableScore {

	for _, testResult := range testResults.TestResults {
		for _, assertionResult := range testResult.AssertionResults {
			for i := 0; i < len(deliverablesScores); i++ {
				if assertionResult.Title == deliverablesScores[i].ScoreCardItemName {
					if assertionResult.Status == "passed" {
						deliverablesScores[i].Pass = true
					}
					break
				}
			}
		}

	}
	return deliverablesScores
}

func buildDeliverableScoresFromPytest(testResults PytestResults, deliverablesScores []DeliverableScore) []DeliverableScore {

	for _, testResult := range testResults.Tests {
		// for _, assertionResult := range testResult.AssertionResults {
			for i := 0; i < len(deliverablesScores); i++ {
				if testResult.Nodeid == "test_residential_controller.py::test_" + deliverablesScores[i].ScoreCardItemName {
					if testResult.Outcome == "passed" {
						deliverablesScores[i].Pass = true
					}
					break
				}
			}
		// }

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
	case "Jest":
		var jestResults JestResults
		if err := json.NewDecoder(strings.NewReader(string(content))).Decode(&jestResults); err != nil {
			panic(err)
		}
		deliverableScores = buildDeliverableScoresFromJest(jestResults, deliverableScores)
	case "Pytest":
		var PytestResults PytestResults
		if err := json.NewDecoder(strings.NewReader(string(content))).Decode(&PytestResults); err != nil {
			panic(err)
		}
		deliverableScores = buildDeliverableScoresFromPytest(PytestResults, deliverableScores)
	}
	return deliverableScores
}
