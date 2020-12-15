package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

//GithubResponse is the data returned by the Github API
type GithubResponse []struct {
	Sha    string `json:"sha"`
	NodeID string `json:"node_id"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			Sha string `json:"sha"`
			URL string `json:"url"`
		} `json:"tree"`
		URL          string `json:"url"`
		CommentCount int    `json:"comment_count"`
		Verification struct {
			Verified  bool        `json:"verified"`
			Reason    string      `json:"reason"`
			Signature interface{} `json:"signature"`
			Payload   interface{} `json:"payload"`
		} `json:"verification"`
	} `json:"commit"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Author      struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Committer struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"committer"`
	Parents []struct {
		Sha     string `json:"sha"`
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"parents"`
}

//CreateIssueRequest contains the JSON to create an issue on github
type CreateIssueRequest struct {
	Title string            `json:"title"`
	Body  CodeClimateIssues `json:"body"`
}

func getLastCommitDate(deliverableScores []DeliverableScore, githubSlug string, deliverableDeadline time.Time) []DeliverableScore {
	var githubResponse GithubResponse
	response, err := http.Get("https://api.github.com/repos/" + githubSlug + "/commits")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		if err := json.NewDecoder(response.Body).Decode(&githubResponse); err != nil {
			panic(err)
		}
	}

	deliveredOnTime := githubResponse[0].Commit.Committer.Date.Before(deliverableDeadline)
	if deliveredOnTime {
		for i := range deliverableScores {
			if deliverableScores[i].ScoreCardItemName == "Delivered on Time" {
				deliverableScores[i].Pass = true
			}
		}
	}
	return deliverableScores
}

func createIssue(githubSlug string, issues CodeClimateIssues) {
	jsonData := CreateIssueRequest{
		Title: time.Now().String() + "_Grading_Request_Code_Analysis",
		Body:  issues,
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, _ := http.NewRequest("POST", "https://api.github.com/repos/"+githubSlug+"/issues", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Token token="+os.Getenv("CODECLIMATETOKEN"))
	request.Header.Set("Accept", "application/vnd.github.v3+json")
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
}
