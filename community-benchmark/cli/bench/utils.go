package bench

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
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

	client := github.NewClient(oAuthClient)

	repoInfo, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return nil, fmt.Errorf("could not get repository info for %s/%s: %s", owner, repo, err)
	}

	pullInfo, issuesInfo, err := collectPRandIssuesInfo(owner, repo, client)
	if err != nil {
		return nil, err
	}

	metrics, err := collectCommunityHealthMetrics(owner, repo, client)
	if err != nil {
		return nil, err
	}

	return &RepositoryInfo{
		Owner:                  owner,
		Repo:                   repo,
		StarNumber:             repoInfo.GetStargazersCount(),
		ForkNumber:             repoInfo.GetForksCount(),
		PullRequestsInfo:       pullInfo,
		IssuesInfo:             issuesInfo,
		CommunityDocumentation: metrics,
		CollectedOn:            time.Now(),
	}, nil
}

func collectPRandIssuesInfo(owner, repo string, client *github.Client) (*PullRequestsInfo, *IssuesInfo, error) {
	var issues []*github.Issue

	pg := 1
	for ; ; pg++ {
		list, _, err := client.Issues.ListByRepo(context.Background(), owner, repo, &github.IssueListByRepoOptions{State: "all", ListOptions: github.ListOptions{PerPage: 100, Page: pg}})
		if err != nil {
			return nil, nil, fmt.Errorf("could not get repository issues for %s/%s: %s", owner, repo, err)
		}

		issues = append(issues, list...)

		if len(list) < 100 {
			break
		}
	}

	openPullRequests := 0
	closedPullRequests := 0
	totalPRComments := 0
	totalPRCommits := 0

	openIssues := 0
	closedIssues := 0
	totalIssueComments := 0
	var wg sync.WaitGroup
	for _, issue := range issues {
		issue := issue
		wg.Add(1)
		go func() {
			if issue.IsPullRequest() {
				if issue.GetState() == "open" {
					openPullRequests++
				}
				if issue.GetState() == "closed" {
					closedPullRequests++
				}
				commits, _, _ := client.PullRequests.ListCommits(context.Background(), owner, repo, issue.GetNumber(), &github.ListOptions{PerPage: 100})
				totalPRCommits += len(commits)

				comments, _, _ := client.PullRequests.ListComments(context.Background(), owner, repo, issue.GetNumber(), &github.PullRequestListCommentsOptions{ListOptions: github.ListOptions{PerPage: 100}})
				totalPRComments += len(comments)
			} else {
				if issue.GetState() == "open" {
					openIssues++
				}
				if issue.GetState() == "closed" {
					closedIssues++
				}
				comments, _, _ := client.Issues.ListComments(context.Background(), owner, repo, issue.GetNumber(), &github.IssueListCommentsOptions{ListOptions: github.ListOptions{PerPage: 100}})
				totalIssueComments += len(comments)
			}

			wg.Done()
		}()
	}
	wg.Wait()

	return &PullRequestsInfo{
			OpenPullRequestNumber:         openPullRequests,
			ClosedPullRequestNumber:       closedPullRequests,
			AverageCommentsPerPullRequest: float64(totalPRComments) / float64(openPullRequests+closedPullRequests),
			AverageCommitsPerPullRequest:  float64(totalPRCommits) / float64(openPullRequests+closedPullRequests),
		}, &IssuesInfo{
			OpenIssuesNumber:        openIssues,
			ClosedIssuesNumber:      closedIssues,
			AverageCommentsPerIssue: float64(totalIssueComments) / float64(openIssues+closedIssues),
		}, nil
}

func collectCommunityHealthMetrics(owner, repo string, client *github.Client) (*CommunityDocumentation, error) {
	metrics, _, err := client.Repositories.GetCommunityHealthMetrics(context.Background(), owner, repo)
	if err != nil {
		return nil, fmt.Errorf("could not get community health metrics for %s/%s: %s", owner, repo, err)
	}

	repoInfo, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return nil, fmt.Errorf("could not get repository info for %s/%s: %s", owner, repo, err)
	}

	return &CommunityDocumentation{
		HealthPercentage:         metrics.GetHealthPercentage(),
		HasCodeOfConduct:         metrics.Files.GetCodeOfConduct() != nil,
		HasContributing:          metrics.Files.GetContributing() != nil,
		HasIssueTemplate:         metrics.Files.GetIssueTemplate() != nil,
		HasPullRequestTemplate:   metrics.Files.GetPullRequestTemplate() != nil,
		HasLicense:               metrics.Files.GetLicense() != nil,
		HasReadme:                metrics.Files.GetReadme() != nil,
		HasContentReportsEnabled: metrics.GetContentReportsEnabled(),
		HasWiki:                  repoInfo.GetHasWiki(),
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
