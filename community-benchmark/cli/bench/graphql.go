package bench

type RepositoryInfoQuery struct {
	Repository struct {
		StargazerCount int
		ForkCount      int
		PullRequests   struct {
			TotalCount int
		}
		Issues struct {
			TotalCount int
		}
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type IssueNode struct {
	Closed   bool
	Comments struct {
		TotalCount int
	}
}

type IssuesInfoQuery struct {
	Repository struct {
		Issues struct {
			PageInfo struct {
				EndCursor string
			}
			Nodes []IssueNode
		} `graphql:"issues(first: 100, after: $endCursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type PullRequestNode struct {
	Closed   bool
	Comments struct {
		TotalCount int
	}
	Commits struct {
		TotalCount int
	}
}

type PullRequestsInfoQuery struct {
	Repository struct {
		PullRequests struct {
			PageInfo struct {
				EndCursor string
			}
			Nodes []PullRequestNode
		} `graphql:"pullRequests(first: 100, after: $endCursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}
