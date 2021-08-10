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

//checkIfRepoIsPrivateResponse contains the IDs of all pending invitations
type checkIfRepoIsPrivateResponse struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
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
	} `json:"owner"`
	HTMLURL          string      `json:"html_url"`
	Description      interface{} `json:"description"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	ForksURL         string      `json:"forks_url"`
	KeysURL          string      `json:"keys_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	TeamsURL         string      `json:"teams_url"`
	HooksURL         string      `json:"hooks_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	EventsURL        string      `json:"events_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BranchesURL      string      `json:"branches_url"`
	TagsURL          string      `json:"tags_url"`
	BlobsURL         string      `json:"blobs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	TreesURL         string      `json:"trees_url"`
	StatusesURL      string      `json:"statuses_url"`
	LanguagesURL     string      `json:"languages_url"`
	StargazersURL    string      `json:"stargazers_url"`
	ContributorsURL  string      `json:"contributors_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	CommitsURL       string      `json:"commits_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	CommentsURL      string      `json:"comments_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	ContentsURL      string      `json:"contents_url"`
	CompareURL       string      `json:"compare_url"`
	MergesURL        string      `json:"merges_url"`
	ArchiveURL       string      `json:"archive_url"`
	DownloadsURL     string      `json:"downloads_url"`
	IssuesURL        string      `json:"issues_url"`
	PullsURL         string      `json:"pulls_url"`
	MilestonesURL    string      `json:"milestones_url"`
	NotificationsURL string      `json:"notifications_url"`
	LabelsURL        string      `json:"labels_url"`
	ReleasesURL      string      `json:"releases_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	PushedAt         time.Time   `json:"pushed_at"`
	GitURL           string      `json:"git_url"`
	SSHURL           string      `json:"ssh_url"`
	CloneURL         string      `json:"clone_url"`
	SvnURL           string      `json:"svn_url"`
	Homepage         interface{} `json:"homepage"`
	Size             int         `json:"size"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Language         string      `json:"language"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasDownloads     bool        `json:"has_downloads"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	ForksCount       int         `json:"forks_count"`
	MirrorURL        interface{} `json:"mirror_url"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	License          interface{} `json:"license"`
	Forks            int         `json:"forks"`
	OpenIssues       int         `json:"open_issues"`
	Watchers         int         `json:"watchers"`
	DefaultBranch    string      `json:"default_branch"`
	Permissions      struct {
		Admin    bool `json:"admin"`
		Maintain bool `json:"maintain"`
		Push     bool `json:"push"`
		Triage   bool `json:"triage"`
		Pull     bool `json:"pull"`
	} `json:"permissions"`
	TempCloneToken   string `json:"temp_clone_token"`
	NetworkCount     int    `json:"network_count"`
	SubscribersCount int    `json:"subscribers_count"`
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

func checkIfGithubUserExists(userName string) bool {
	request, _ := http.NewRequest("GET", "https://api.github.com/users/"+userName, nil)
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{}
	response, _ := client.Do(request)
	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return true
	} else {
		return false
	}
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

func checkIfRepoIsPrivate(repositoryURL string) bool {
	githubSlug := strings.Replace(repositoryURL[strings.LastIndex(repositoryURL, ":")+1:], ".git", "", -1)
	request, _ := http.NewRequest("GET", "https://api.github.com/repos/"+githubSlug, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "token "+os.Getenv("GITHUBTOKEN"))
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{}
	response, err := client.Do(request)
	var jsonResponseBody checkIfRepoIsPrivateResponse
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return false
	} else if response.StatusCode == 404 {
		return true
	} else {

		if err := json.NewDecoder(response.Body).Decode(&jsonResponseBody); err != nil {
			panic(err)
		}

		return jsonResponseBody.Private

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
