package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

//GetRepositoryInvitationsResponse contains the IDs of all pending invitations
type GetRepositoryInvitationsResponse []struct {
	ID int `json:"id"`
}

//ForkRepositoryResponse contains the id and the name of the forked repository
type ForkRepositoryResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}

//GetLastCommitResponse is the data returned by the Github API
type GetLastCommitResponse []struct {
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
	Title string `json:"title"`
	Body  string `json:"body"`
}

func checkGithubAccess(repositoryURL string) bool {
	githubSlug := strings.Replace(repositoryURL[strings.LastIndex(repositoryURL, ":")+1:], ".git", "", -1)
	request, _ := http.NewRequest("GET", "https://api.github.com/repos/"+githubSlug, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "token "+os.Getenv("GITHUBTOKEN"))
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{}
	response, _ := client.Do(request)
	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return true
	} else {
		return false
	}
}

func getLastCommitDate(githubSlug string) time.Time {
	var lastCommit GetLastCommitResponse
	request, _ := http.NewRequest("GET", "https://api.github.com/repos/"+githubSlug+"/commits", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "token "+os.Getenv("GITHUBTOKEN"))
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		if err := json.NewDecoder(response.Body).Decode(&lastCommit); err != nil {
			panic(err)
		}
	}
	lastCommitDate := lastCommit[0].Commit.Committer.Date

	return lastCommitDate
}

func buildGithubIssueBody(issues CodeClimateIssues) string {
	var sb strings.Builder
	for _, issue := range issues.Data {
		sb.WriteString("Issue: " + issue.Attributes.CheckName + " \n")
		sb.WriteString("\t - File: " + issue.Attributes.ConstantName + " \n")
		sb.WriteString("\t - Description: " + issue.Attributes.Description + " \n")
		sb.WriteString(fmt.Sprintf("\t - Location: From line %d to %d. \n", +issue.Attributes.Location.StartLine, issue.Attributes.Location.EndLine))
		sb.WriteString("\t - Severity: " + issue.Attributes.Severity + " \n \n")
	}

	return sb.String()
}

func createIssue(githubSlug string, issues CodeClimateIssues) {
	issueBody := buildGithubIssueBody(issues)
	jsonData := CreateIssueRequest{
		Title: time.Now().Format("2006-01-02T15:04:05.999999") + "_Grading_Request_Code_Analysis",
		Body:  issueBody,
	}
	jsonValue, _ := json.Marshal(jsonData)
	request, _ := http.NewRequest("POST", "https://api.github.com/repos/"+githubSlug+"/issues", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "token "+os.Getenv("GITHUBTOKEN"))
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{}
	response, err := client.Do(request)
	var jsonResponseBody AddRepositoryResponse
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		if err := json.NewDecoder(response.Body).Decode(&jsonResponseBody); err != nil {
			panic(err)
		}

	}
}

func checkRepositoryInvitations() []int {
	var repositoryInvitationsIDs GetRepositoryInvitationsResponse
	request, _ := http.NewRequest("GET", "https://api.github.com/user/repository_invitations", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "token "+os.Getenv("GITHUBTOKEN"))
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		if err := json.NewDecoder(response.Body).Decode(&repositoryInvitationsIDs); err != nil {
			panic(err)
		}

	}
	var IDList []int
	for _, invitation := range repositoryInvitationsIDs {
		IDList = append(IDList, invitation.ID)
	}

	return IDList
}

func acceptRepositoryInvitations(IDList []int) {
	for _, invitationID := range IDList {
		request, _ := http.NewRequest("PATCH", "https://api.github.com/user/repository_invitations/"+fmt.Sprint(invitationID), nil)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "token "+os.Getenv("GITHUBTOKEN"))
		request.Header.Set("Accept", "application/vnd.github.v3+json")
		client := &http.Client{}
		_, err := client.Do(request)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}
	}
}

func forkRepository(githubSlug string) string {
	jsonData := struct {
		Organization string `json:"organization"`
	}{Organization: "Codeboxx-Students-Projects"}

	jsonValue, _ := json.Marshal(jsonData)
	request, _ := http.NewRequest("POST", "https://api.github.com/repos/"+githubSlug+"/forks", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "token "+os.Getenv("GITHUBTOKEN"))
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{}
	var forkRepositoryResponse ForkRepositoryResponse
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		if err := json.NewDecoder(response.Body).Decode(&forkRepositoryResponse); err != nil {
			panic(err)
		}
	}

	return forkRepositoryResponse.FullName
}
