package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

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

	IDList := checkRepositoryInvitations()
	if len(IDList) > 0 {
		acceptRepositoryInvitations(IDList)
	}

	githubSlug := strings.Replace(gradingRequest.RepositoryURL[strings.LastIndex(gradingRequest.RepositoryURL, ":")+1:], ".git", "", -1)
	deliverableScores := docker(gradingRequest)
	deliverableScores = checkRespectOfDeadline(deliverableScores, githubSlug, gradingRequest)
	forkedRepoName := forkRepository(githubSlug)
	issues := codeClimate(forkedRepoName, deliverableScores)
	createIssue(githubSlug, issues)

	gradingResponse := GradingResponse{
		JobID:             gradingRequest.JobID,
		DeliverableID:     gradingRequest.DeliverableID,
		DeliverableScores: deliverableScores,
		Issues:            issues,
	}

	return gradingResponse

}

func checkRespectOfDeadline(deliverableScores []DeliverableScore, githubSlug string, gradingRequest GradingRequest) []DeliverableScore {
	lastCommitDate := getLastCommitDate(deliverableScores, githubSlug)
	deliveredOnTime := lastCommitDate.Before(gradingRequest.DeliverableDeadline)
	if deliveredOnTime {
		for i := range deliverableScores {
			if deliverableScores[i].ScoreCardItemName == "Delivered on Time" {
				deliverableScores[i].Pass = true
			}
		}
	}
	return deliverableScores
}
