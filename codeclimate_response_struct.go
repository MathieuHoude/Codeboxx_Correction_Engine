package main

import "time"

//AddPublicRepositoryResponse contains the data returned by codeclimate
type AddPublicRepositoryResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			AnalysisVersion       int         `json:"analysis_version"`
			BadgeToken            string      `json:"badge_token"`
			Branch                string      `json:"branch"`
			CreatedAt             time.Time   `json:"created_at"`
			DelegatedConfigRepoID string      `json:"delegated_config_repo_id"`
			DiffCoverageEnforced  bool        `json:"diff_coverage_enforced"`
			DiffCoverageThreshold int         `json:"diff_coverage_threshold"`
			EnableNotifications   bool        `json:"enable_notifications"`
			GithubSlug            string      `json:"github_slug"`
			HumanName             string      `json:"human_name"`
			LastActivityAt        time.Time   `json:"last_activity_at"`
			TestReporterID        string      `json:"test_reporter_id"`
			TotalCoverageEnforced bool        `json:"total_coverage_enforced"`
			VcsDatabaseID         string      `json:"vcs_database_id"`
			VcsHost               string      `json:"vcs_host"`
			Score                 interface{} `json:"score"`
		} `json:"attributes"`
		Relationships struct {
			LatestDefaultBranchSnapshot struct {
				Data interface{} `json:"data"`
			} `json:"latest_default_branch_snapshot"`
			LatestDefaultBranchTestReport struct {
				Data interface{} `json:"data"`
			} `json:"latest_default_branch_test_report"`
			Account struct {
				Data interface{} `json:"data"`
			} `json:"account"`
		} `json:"relationships"`
		Links struct {
			Self                 string `json:"self"`
			Services             string `json:"services"`
			WebCoverage          string `json:"web_coverage"`
			WebIssues            string `json:"web_issues"`
			MaintainabilityBadge string `json:"maintainability_badge"`
			TestCoverageBadge    string `json:"test_coverage_badge"`
		} `json:"links"`
		Meta struct {
			Permissions struct {
				Admin bool `json:"admin"`
			} `json:"permissions"`
		} `json:"meta"`
	} `json:"data"`
}

//GetRepositoryResponse contains the data returned by getRepository
type GetRepositoryResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			AnalysisVersion       int         `json:"analysis_version"`
			BadgeToken            string      `json:"badge_token"`
			Branch                string      `json:"branch"`
			CreatedAt             time.Time   `json:"created_at"`
			DelegatedConfigRepoID string      `json:"delegated_config_repo_id"`
			DiffCoverageEnforced  bool        `json:"diff_coverage_enforced"`
			DiffCoverageThreshold int         `json:"diff_coverage_threshold"`
			EnableNotifications   bool        `json:"enable_notifications"`
			GithubSlug            string      `json:"github_slug"`
			HumanName             string      `json:"human_name"`
			LastActivityAt        time.Time   `json:"last_activity_at"`
			TestReporterID        interface{} `json:"test_reporter_id"`
			TotalCoverageEnforced bool        `json:"total_coverage_enforced"`
			VcsDatabaseID         string      `json:"vcs_database_id"`
			VcsHost               string      `json:"vcs_host"`
			Score                 interface{} `json:"score"`
		} `json:"attributes"`
		Relationships struct {
			LatestDefaultBranchSnapshot struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"latest_default_branch_snapshot"`
			LatestDefaultBranchTestReport struct {
				Data interface{} `json:"data"`
			} `json:"latest_default_branch_test_report"`
			Account struct {
				Data interface{} `json:"data"`
			} `json:"account"`
		} `json:"relationships"`
		Links struct {
			Self                 string `json:"self"`
			Services             string `json:"services"`
			WebCoverage          string `json:"web_coverage"`
			WebIssues            string `json:"web_issues"`
			MaintainabilityBadge string `json:"maintainability_badge"`
			TestCoverageBadge    string `json:"test_coverage_badge"`
		} `json:"links"`
		Meta struct {
			Permissions struct {
				Admin bool `json:"admin"`
			} `json:"permissions"`
		} `json:"meta"`
	} `json:"data"`
}

// CodeClimateIssues contains the data returned by codeclimate for the getIssues endpoint
type CodeClimateIssues struct {
	Data []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Categories   []string `json:"categories"`
			CheckName    string   `json:"check_name"`
			ConstantName string   `json:"constant_name"`
			Content      struct {
				Body string `json:"body"`
			} `json:"content"`
			Description string `json:"description"`
			Location    struct {
				Path      string `json:"path"`
				EndLine   int    `json:"end_line"`
				StartLine int    `json:"start_line"`
			} `json:"location"`
			OtherLocations    []interface{} `json:"other_locations"`
			RemediationPoints int           `json:"remediation_points"`
			Severity          string        `json:"severity"`
		} `json:"attributes"`
	} `json:"data"`
}
