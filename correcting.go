package main

// CorrectionRequest contains the necessary elements to correct a project
type CorrectionRequest struct {
	JobID                   uint   `json:"JobID"`
	DeliverableID           uint   `json:"DeliverableID"`
	UnixDeliverableDeadline uint64 `json:"UnixDeliverableDeadline"`
	GithubHandle            string `json:"GithubHandle"`
	RepositoryURL           string `json:"RepositoryURL"`
	DockerImageName         string `json:"DockerImageName"`
}

func startCorrecting(correctionRequest CorrectionRequest) {

	IDList := checkRepositoryInvitations()
	if len(IDList) > 0 {
		acceptRepositoryInvitations(IDList)
	}

	docker(correctionRequest)
}
