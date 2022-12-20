package bench

import "time"

type RepositoryInfo struct {
	Owner                  string                  `json:"owner"`
	Repo                   string                  `json:"repo"`
	StarNumber             int                     `json:"star_number"`
	ForkNumber             int                     `json:"fork_number"`
	PullRequestsInfo       *PullRequestsInfo       `json:"pull_requests_info"`
	IssuesInfo             *IssuesInfo             `json:"issues_info"`
	CommunityDocumentation *CommunityDocumentation `json:"community_documentation"`
	CollectedOn            time.Time               `json:"collected_on"`
}

type PullRequestsInfo struct {
	OpenPullRequestNumber         int     `json:"open_pull_request_number"`
	ClosedPullRequestNumber       int     `json:"closed_pull_request_number"`
	AverageCommentsPerPullRequest float64 `json:"average_comments_per_pull_request"`
	AverageCommitsPerPullRequest  float64 `json:"average_commits_per_pull_request"`
}

type IssuesInfo struct {
	OpenIssuesNumber        int     `json:"open_issues_number"`
	ClosedIssuesNumber      int     `json:"closed_issues_number"`
	AverageCommentsPerIssue float64 `json:"average_comments_per_issue"`
}

type CollaboratorsInfo struct {
}

type CommunityDocumentation struct {
	HealthPercentage         int  `json:"health_percentage"`
	HasCodeOfConduct         bool `json:"has_code_of_conduct"`
	HasContributing          bool `json:"has_contributing"`
	HasIssueTemplate         bool `json:"has_issue_template"`
	HasPullRequestTemplate   bool `json:"has_pull_request_template"`
	HasLicense               bool `json:"has_license"`
	HasReadme                bool `json:"has_readme"`
	HasContentReportsEnabled bool `json:"has_content_reports_enabled"`
	HasWiki                  bool `json:"has_wiki"`
	// HasDiscussions           bool `json:"has_discussions"`
}
