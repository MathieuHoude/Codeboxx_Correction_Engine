package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//DeliverableScores contains the list of scores for a given deliverable
type DeliverableScores struct {
	DeliverableID     int
	DeliverableScores []struct {
		ID   int
		Pass bool
	}
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
	jsonData := []byte(`{"JobID": "` + fmt.Sprint(jobID) + `", NewStatus: "` + newStatus + `"}`)
	request, _ := http.NewRequest("POST", os.Getenv("JOBUPDATEENDPOINT"), bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	failOnError(err, "Failed to update job status")
	defer response.Body.Close()
}

func startGrading(gradingRequest GradingRequest) GradingResponse {
	var gradingResponse GradingResponse
	switch gradingRequest.ProjectName {
	case "ruby-residential-controller":
		gradingResponse = rubyResidentialControllerCorrection(gradingRequest)

	}
	return gradingResponse

}

func buildDeliverableScoresFromRspec(testResults RspecResults) DeliverableScores {
	var scores DeliverableScores
	return scores
}

func rubyResidentialControllerCorrection(gradingRequest GradingRequest) GradingResponse {
	testResults := docker(gradingRequest)
	issues := codeClimate(gradingRequest.GithubHandle, "Rocket_Elevators_Controllers")

	scores := buildDeliverableScoresFromRspec(testResults)

	gradingResponse := GradingResponse{JobID: gradingRequest.JobID, DeliverableScores: scores, Issues: issues}

	return gradingResponse
}
