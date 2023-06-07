package community

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func extractRepo(s string) (string, string, error) {
	ownerPattern := "[a-zA-Z0-9]([a-zA-Z0-9]|-[a-zA-Z0-9]){0,38}"
	repoPattern := ownerPattern
	validator := regexp.MustCompile("^(https://)?github.com/" + ownerPattern + "/" + repoPattern + "/?$")
	if !validator.MatchString(s) {
		return "", "", fmt.Errorf("provided argument %s is not a valid github repository link", s)
	}

	trimmer := regexp.MustCompile("^(https://)?github.com/")
	prefixTrimmed := trimmer.ReplaceAllString(s, "")

	splits := strings.Split(prefixTrimmed, "/")
	if len(splits[0]) > 39 || len(splits[1]) > 39 {
		return "", "", fmt.Errorf("provided argument %s is not a valid github repository link", s)
	}

	return splits[0], splits[1], nil
}

func collectRepositoryInfo(owner, repo string) (*RepositoryInfo, error) {
	var oAuthClient *http.Client = nil
	accessToken, ok := os.LookupEnv("GITHUB_TOKEN")
	if ok {
		oAuthClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		))
	}

	client := githubv4.NewClient(oAuthClient)

	repoInfo := RepositoryInfoQuery{}
	err := client.Query(context.Background(), &repoInfo, map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
	})
	if err != nil {
		return nil, err
	}

	issuesInfo, err := collectIssuesInfo(owner, repo, repoInfo.Repository.Issues.TotalCount, client)
	if err != nil {
		return nil, err
	}

	pullRequestsInfo, err := collectPullRequestsInfo(owner, repo, repoInfo.Repository.PullRequests.TotalCount, client)
	if err != nil {
		return nil, err
	}

	//metrics, err := collectCommunityHealthMetrics(owner, repo, client)
	//if err != nil {
	//	return nil, err
	//}

	return &RepositoryInfo{
		Owner:            owner,
		Repo:             repo,
		StarNumber:       repoInfo.Repository.StargazerCount,
		ForkNumber:       repoInfo.Repository.ForkCount,
		PullRequestsInfo: pullRequestsInfo,
		IssuesInfo:       issuesInfo,
		//CommunityDocumentation: metrics,
		CollectedOn: time.Now(),
	}, nil
}

func collectIssuesInfo(owner, repo string, issueCount int, client *githubv4.Client) (*IssuesInfo, error) {
	var issues []IssueNode

	endCursor := (*githubv4.String)(nil)
	for index := 0; index < issueCount; index += 100 {
		issuesInfo := IssuesInfoQuery{}
		err := client.Query(context.Background(), &issuesInfo, map[string]interface{}{
			"owner":     githubv4.String(owner),
			"name":      githubv4.String(repo),
			"endCursor": endCursor,
		})
		if err != nil {
			return nil, err
		}

		endCursor = githubv4.NewString(githubv4.String(issuesInfo.Repository.Issues.PageInfo.EndCursor))

		issues = append(issues, issuesInfo.Repository.Issues.Nodes...)
	}

	openIssues := 0
	closedIssues := 0
	totalComments := 0
	for _, issue := range issues {
		if issue.Closed {
			closedIssues++
		} else {
			openIssues++
		}
		totalComments += issue.Comments.TotalCount
	}

	return &IssuesInfo{
		OpenIssuesNumber:        openIssues,
		ClosedIssuesNumber:      closedIssues,
		AverageCommentsPerIssue: float64(totalComments) / float64(issueCount),
	}, nil
}

func collectPullRequestsInfo(owner, repo string, pullRequestCount int, client *githubv4.Client) (*PullRequestsInfo, error) {
	var pullRequests []PullRequestNode

	endCursor := (*githubv4.String)(nil)
	for index := 0; index < pullRequestCount; index += 100 {
		pullRequestInfo := PullRequestsInfoQuery{}
		err := client.Query(context.Background(), &pullRequestInfo, map[string]interface{}{
			"owner":     githubv4.String(owner),
			"name":      githubv4.String(repo),
			"endCursor": endCursor,
		})
		if err != nil {
			return nil, err
		}

		endCursor = githubv4.NewString(githubv4.String(pullRequestInfo.Repository.PullRequests.PageInfo.EndCursor))

		pullRequests = append(pullRequests, pullRequestInfo.Repository.PullRequests.Nodes...)
	}

	openPullRequests := 0
	closedPullRequests := 0
	totalComments := 0
	totalCommits := 0
	for _, pullRequest := range pullRequests {
		if pullRequest.Closed {
			closedPullRequests++
		} else {
			openPullRequests++
		}
		totalComments += pullRequest.Comments.TotalCount
		totalCommits += pullRequest.Commits.TotalCount
	}

	return &PullRequestsInfo{
		OpenPullRequestsNumber:        openPullRequests,
		ClosedPullRequestsNumber:      closedPullRequests,
		AverageCommentsPerPullRequest: float64(totalComments) / float64(pullRequestCount),
		AverageCommitsPerPullRequest:  float64(totalCommits) / float64(pullRequestCount),
	}, nil
}

func saveOutput(c *cli.Context, info *RepositoryInfo) error {
	if c.Bool("mysql-export") {
		connection, err := handleMySQLConnection()
		if err != nil {
			return err
		}

		err = saveToDb(connection, info)
		if err != nil {
			return err
		}

		err = connection.Close()
		if err != nil {
			return err
		}

		log.Println("Contents successfully saved to MySQL database")

		return nil
	}

	marshal, err := json.MarshalIndent(info, "", "	")
	if err != nil {
		return err
	}

	fmt.Println(string(marshal))

	return nil
}
