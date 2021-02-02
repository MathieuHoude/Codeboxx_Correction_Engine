package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

//JSONRequestBody contains the necessary data for the API call
type JSONRequestBody struct {
	Data Data `json:"data"`
}

//Attributes is the url of the repo
type Attributes struct {
	URL string `json:"url"`
}

//Data is the two main elements of the request body
type Data struct {
	ID         string
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

func codeClimateHTTPRequest(method, URL string, body []byte) (*http.Response, error) {
	request, _ := http.NewRequest(method, URL, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Token token="+os.Getenv("CODECLIMATETOKEN"))
	request.Header.Set("Accept", "application/vnd.api+json")
	client := &http.Client{}
	response, err := client.Do(request)

	return response, err
}

func addPublicRepository(githubSlug string) AddRepositoryResponse {
	jsonData := JSONRequestBody{
		Data{
			Type: "repos",
			Attributes: Attributes{
				URL: "https://github.com/" + githubSlug,
			},
		},
	}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := codeClimateHTTPRequest("POST", "https://api.codeclimate.com/v1/github/repos", jsonValue)
	var jsonResponseBody AddRepositoryResponse
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		if err := json.NewDecoder(response.Body).Decode(&jsonResponseBody); err != nil {
			panic(err)
		}

	}
	return jsonResponseBody
}

func addPrivateRepository(githubSlug string) AddRepositoryResponse {
	jsonData := JSONRequestBody{
		Data{
			Type: "repos",
			Attributes: Attributes{
				URL: "https://github.com/" + githubSlug,
			},
		},
	}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := codeClimateHTTPRequest("POST", "https://api.codeclimate.com/v1/orgs/"+os.Getenv("CODECLIMATEORGANIZATIONID")+"/repos", jsonValue)
	var jsonResponseBody AddRepositoryResponse
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		if err := json.NewDecoder(response.Body).Decode(&jsonResponseBody); err != nil {
			panic(err)
		}

	}
	return jsonResponseBody
}

//Needs refactoring
func getPrivateRepository(githubSlug string) GetRepositoryResponse {
	var jsonResponseBody GetRepositoryResponse
	for jsonResponseBody.Data == nil || jsonResponseBody.Data[0].Relationships.LatestDefaultBranchSnapshot.Data.ID == "" {
		time.Sleep(2 * time.Second)
		response, err := codeClimateHTTPRequest("GET", "https://api.codeclimate.com/v1/repos?github_slug="+githubSlug, nil)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			if err := json.NewDecoder(response.Body).Decode(&jsonResponseBody); err != nil {
				panic(err)
			}
		}
	}
	return jsonResponseBody
}

func getPrivateRepositoryIssues(getRepositoryResponse GetRepositoryResponse) CodeClimateIssues {
	repoID := getRepositoryResponse.Data[0].ID
	snapshotID := getRepositoryResponse.Data[0].Relationships.LatestDefaultBranchSnapshot.Data.ID
	var codeClimateIssues CodeClimateIssues
	response, err := codeClimateHTTPRequest("GET", "https://api.codeclimate.com/v1/repos/"+repoID+"/snapshots/"+snapshotID+"/issues", nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		if err := json.NewDecoder(response.Body).Decode(&codeClimateIssues); err != nil {
			panic(err)
		}
	}

	return codeClimateIssues
}

func codeClimate(githubSlug string, deliverableScores []DeliverableScore) CodeClimateIssues {
	addPrivateRepository(githubSlug)
	getRepositoryResponse := getPrivateRepository(githubSlug)
	codeClimateIssues := getPrivateRepositoryIssues(getRepositoryResponse)

	if len(codeClimateIssues.Data) < 3 {
		for i := range deliverableScores {
			if deliverableScores[i].ScoreCardItemName == "The program is well coded" {
				deliverableScores[i].Pass = true
			}
		}
	}

	return codeClimateIssues
}
