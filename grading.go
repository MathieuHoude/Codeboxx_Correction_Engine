package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// GradingRequest contains the necessary elements to grade a project
type GradingRequest struct {
	JobID                   uint         `json:"JobID"`
	DeliverableID           uint         `json:"DeliverableID"`
	UnixDeliverableDeadline uint64       `json:"UnixDeliverableDeadline"`
	RepositoryURL           string       `json:"RepositoryURL"`
	TestResults             []TestResult `json:"TestResults"`
}

//GradingResponse contains the informations to send back to the requester
type GradingResponse struct {
	DeliverableID     uint
	DeliverableScores []DeliverableScore
	Issues            CodeClimateIssues
}

func sendBackResults(jobID uint, gradingResponse GradingResponse) {
	jsonData, _ := json.Marshal(gradingResponse)
	fmt.Println(string(jsonData))
	request, _ := http.NewRequest("POST", os.Getenv("DELIVERABLESCORESUPDATEENDPOINT"), bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	failOnError(err, "Failed to send back grading results")
	defer response.Body.Close()

	fmt.Println("Finished job #" + fmt.Sprint(jobID))
}

//newStatus should be max 20 caracters long
func updateJobStatus(jobID uint, newStatus string) {
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
	deliverableScores := buildDeliverableScores(gradingRequest.TestResults)
	githubSlug := strings.Replace(gradingRequest.RepositoryURL[strings.LastIndex(gradingRequest.RepositoryURL, ":")+1:], ".git", "", -1)
	deliveredOnTimeScore := checkRespectOfDeadline(githubSlug, gradingRequest.UnixDeliverableDeadline)
	deliverableScores = append(deliverableScores, deliveredOnTimeScore)
	// forkedRepoName := forkRepository(githubSlug)
	// issues := codeClimate(forkedRepoName, deliverableScores)
	// createIssue(githubSlug, issues)

	gradingResponse := GradingResponse{
		DeliverableID:     gradingRequest.DeliverableID,
		DeliverableScores: deliverableScores,
	}

	return gradingResponse
}

func checkRespectOfDeadline(githubSlug string, UnixDeliverableDeadlineUnix uint64) DeliverableScore {
	deliveredOnTimeScore := DeliverableScore{"Delivered on Time", false}
	lastCommitDate := getLastCommitDate(githubSlug)
	i, err := strconv.ParseInt(fmt.Sprint(UnixDeliverableDeadlineUnix), 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	deliveredOnTime := lastCommitDate.Before(tm)

	if deliveredOnTime {
		deliveredOnTimeScore.Pass = true
	}
	return deliveredOnTimeScore
}
