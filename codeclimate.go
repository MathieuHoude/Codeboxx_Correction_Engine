package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func addPublicRepository(githubSlug string) AddPublicRepositoryResponse {
	jsonData := JSONRequestBody{
		Data{
			Type: "repos",
			Attributes: Attributes{
				URL: "https://github.com/" + githubSlug,
			},
		},
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, _ := http.NewRequest("POST", "https://api.codeclimate.com/v1/github/repos", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Token token=85ab81759fb277a547e26acfc8b29b461030d449")
	request.Header.Set("Accept", "application/vnd.api+json")
	client := &http.Client{}
	response, err := client.Do(request)
	var jsonResponseBody AddPublicRepositoryResponse
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
func getRepository(githubSlug string) GetRepositoryResponse {
	var jsonResponseBody GetRepositoryResponse
	for jsonResponseBody.Data == nil || jsonResponseBody.Data[0].Relationships.LatestDefaultBranchSnapshot.Data.ID == "" {
		time.Sleep(2 * time.Second)
		response, err := http.Get("https://api.codeclimate.com/v1/repos?github_slug=" + githubSlug)
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

func getIssues(getRepositoryResponse GetRepositoryResponse) CodeClimateIssues {
	repoID := getRepositoryResponse.Data[0].ID
	snapshotID := getRepositoryResponse.Data[0].Relationships.LatestDefaultBranchSnapshot.Data.ID
	var codeClimateIssues CodeClimateIssues
	response, err := http.Get("https://api.codeclimate.com/v1/repos/" + repoID + "/snapshots/" + snapshotID + "/issues")

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		if err := json.NewDecoder(response.Body).Decode(&codeClimateIssues); err != nil {
			panic(err)
		}
	}

	return codeClimateIssues
}

func codeClimate(githubSlug string) CodeClimateIssues {
	addPublicRepository(githubSlug)
	getRepositoryResponse := getRepository(githubSlug)
	codeClimateIssues := getIssues(getRepositoryResponse)
	return codeClimateIssues
}
