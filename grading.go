package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

//DeliverableScore represents a deliverable criteria
type DeliverableScore struct {
	ID                int
	ScoreCardItemName string
	Pass              bool
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

func sendBackResults(gradingResponse GradingResponse) {
	jsonData, _ := json.Marshal(gradingResponse)
	request, _ := http.NewRequest("POST", os.Getenv("DELIVERABLESCORESUPDATEENDPOINT"), bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	failOnError(err, "Failed to send back grading results")
	defer response.Body.Close()

	fmt.Println("Finished job #" + fmt.Sprint(gradingResponse.JobID))
}

func updateJobStatus(jobID int, newStatus string) {
	jsonData := []byte(`{"GradingJobsID": "` + fmt.Sprint(jobID) + `", Results: "` + newStatus + `"}`)
	request, _ := http.NewRequest("POST", os.Getenv("JOBUPDATEENDPOINT"), bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	failOnError(err, "Failed to update job status")
	defer response.Body.Close()
}

func startGrading(gradingRequest GradingRequest) GradingResponse {
	deliverableScores := docker(gradingRequest)
	issues := codeClimate(gradingRequest.RepositoryURL)

	gradingResponse := GradingResponse{JobID: gradingRequest.JobID, DeliverableScores: deliverableScores, Issues: issues}

	return gradingResponse

}

func buildDeliverableScoresFromRspec(testResults RspecResults, deliverablesScores []DeliverableScore) []DeliverableScore {

	for _, testResult := range testResults.Examples {
		for _, deliverableScore := range deliverablesScores {
			if testResult.Description == deliverableScore.ScoreCardItemName {
				if testResult.Status == "passed" {
					deliverableScore.Pass = false
				}
				break
			}
		}

	}
	return deliverablesScores
}
