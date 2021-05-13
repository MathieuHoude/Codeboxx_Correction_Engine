package main

//DeliverableScore represents a deliverable criteria
type DeliverableScore struct {
	ScoreCardItemName string
	Pass              bool
}

// func buildDeliverableScoresFromRspec(testResults RspecResults, deliverablesScores []DeliverableScore) []DeliverableScore {

// 	for _, testResult := range testResults.Examples {
// 		for i := 0; i < len(deliverablesScores); i++ {
// 			if testResult.Description == deliverablesScores[i].ScoreCardItemName {
// 				if testResult.Status == "passed" {
// 					deliverablesScores[i].Pass = true
// 				}
// 				break
// 			}
// 		}

// 	}
// 	return deliverablesScores
// }

func buildDeliverableScores(testResults []TestResult) []DeliverableScore {
	var deliverableScores []DeliverableScore
	for _, testResult := range testResults {
		deliverableScore := DeliverableScore{testResult.ScoreCardItemName, false}
		if testResult.Status == "passed" {
			deliverableScore.Pass = true
		}
		deliverableScores = append(deliverableScores, deliverableScore)
	}
	return deliverableScores
}

// func buildDeliverableScores(gradingRequest GradingRequest) []DeliverableScore {
// 	var deliverableScores []DeliverableScore

// 	switch gradingRequest.TestingTool {
// 	case "Rspec":
// 		var rspecResults RspecResults
// 		if err := json.NewDecoder(strings.NewReader(string(content))).Decode(&rspecResults); err != nil {
// 			panic(err)
// 		}
// 		deliverableScores = buildDeliverableScoresFromRspec(rspecResults, deliverableScores)
// 	case "Jest":
// 		// var testResults TestResults
// 		// if err := json.NewDecoder(strings.NewReader(string(content))).Decode(&jestResults); err != nil {
// 		// 	panic(err)
// 		// }
// 		deliverableScores = buildDeliverableScoresFromJest(gradingRequest.TestResults)
// 	case "Pytest":
// 		var PytestResults PytestResults
// 		if err := json.NewDecoder(strings.NewReader(string(content))).Decode(&PytestResults); err != nil {
// 			panic(err)
// 		}
// 		deliverableScores = buildDeliverableScoresFromPytest(PytestResults, deliverableScores)
// 	}
// 	return deliverableScores
// }
