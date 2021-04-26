package main

// CorrectionRequest contains the necessary elements to correct a project
type CorrectionRequest struct {
	JobID               uint   `json:"JobID"`
	DeliverableID       uint   `json:"DeliverableID"`
	DeliverableDeadline uint64 `json:"DeliverableDeadline"`
	GithubHandle        string `json:"GithubHandle"`
	RepositoryURL       string `json:"RepositoryURL"`
	DockerImageName     string `json:"DockerImageName"`
}

func startCorrecting(correctionRequest CorrectionRequest) {

	IDList := checkRepositoryInvitations()
	if len(IDList) > 0 {
		acceptRepositoryInvitations(IDList)
	}

	docker(correctionRequest)

	// githubSlug := strings.Replace(gradingRequest.RepositoryURL[strings.LastIndex(gradingRequest.RepositoryURL, ":")+1:], ".git", "", -1)
	// deliverableScores := docker(gradingRequest)
	// deliverableScores := checkRespectOfDeadline(gradingRequest.DeliverableScores, githubSlug, gradingRequest)
	// forkedRepoName := forkRepository(githubSlug)
	// issues := codeClimate(forkedRepoName, deliverableScores)
	// createIssue(githubSlug, issues)

	// gradingResponse := GradingResponse{
	// 	JobID:             gradingRequest.JobID,
	// 	DeliverableID:     gradingRequest.DeliverableID,
	// 	DeliverableScores: deliverableScores,
	// 	Issues:            issues,
	// }

	// return gradingResponse

}
