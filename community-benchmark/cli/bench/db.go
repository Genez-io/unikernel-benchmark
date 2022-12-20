package bench

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"time"
)

func getEnvVariable(name string, private bool) (string, error) {
	envVar, ok := os.LookupEnv(name)
	if !ok {
		fmt.Println(name + " environment variable is not set. Please input a value:")
		_, err := fmt.Scanln(&envVar)
		if err != nil {
			return "", err
		}

		return envVar, nil
	}

	if private {
		fmt.Println("Using " + name + " environment variable: " + strings.Repeat("*", len(envVar)))
	} else {
		fmt.Println("Using " + name + " environment variable: " + envVar)
	}
	return envVar, nil
}

func handleMySQLConnection() (*sql.DB, error) {
	databaseAddr, err := getEnvVariable("MYSQL_ADDRESS", false)
	if err != nil {
		return nil, err
	}
	databasePort, err := getEnvVariable("MYSQL_PORT", false)
	if err != nil {
		return nil, err
	}
	user, err := getEnvVariable("MYSQL_USER", false)
	if err != nil {
		return nil, err
	}
	password, err := getEnvVariable("MYSQL_PASSWORD", true)
	if err != nil {
		return nil, err
	}
	databaseName, err := getEnvVariable("MYSQL_DATABASE", false)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, databaseAddr, databasePort, databaseName))
	for try := 0; try < 10; try++ {
		err := conn.Ping()
		if err != nil {
			if try == 9 {
				return nil, err
			}
			log.Println("Failed to connect to database. Retrying...")
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	return conn, nil
}

func saveToDb(connection *sql.DB, info *RepositoryInfo) error {
	tx, err := connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		`INSERT INTO repository (owner, repo, star_number, fork_number) VALUES (?, ?, ?, ?)`,
		info.Owner,
		info.Repo,
		info.StarNumber,
		info.ForkNumber)
	if err != nil {
		return err
	}
	repositoryId, _ := result.LastInsertId()

	result, err = tx.Exec(
		`INSERT INTO pull_requests (repository_id, open_pull_requests_number, closed_pull_requests_number, 
			average_comments_per_pull_request, average_commits_per_pull_request) VALUES (?, ?, ?, ?, ?)`,
		repositoryId,
		info.PullRequestsInfo.OpenPullRequestNumber,
		info.PullRequestsInfo.ClosedPullRequestNumber,
		info.PullRequestsInfo.AverageCommentsPerPullRequest,
		info.PullRequestsInfo.AverageCommitsPerPullRequest)
	if err != nil {
		return err
	}

	result, err = tx.Exec(
		`INSERT INTO issues (repository_id, open_issues_number, closed_issues_number,
        	average_comments_per_issue) VALUES (?, ?, ?, ?)`,
		repositoryId,
		info.IssuesInfo.OpenIssuesNumber,
		info.IssuesInfo.ClosedIssuesNumber,
		info.IssuesInfo.AverageCommentsPerIssue)
	if err != nil {
		return err
	}

	result, err = tx.Exec(
		`INSERT INTO community_documents (repository_id, health_percentage, has_code_of_conduct,
			has_contributing, has_issue_template, has_pull_request_template, has_license, has_readme, 
           	has_content_reports_enabled, has_wiki) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		repositoryId,
		info.CommunityDocumentation.HealthPercentage,
		info.CommunityDocumentation.HasCodeOfConduct,
		info.CommunityDocumentation.HasContributing,
		info.CommunityDocumentation.HasIssueTemplate,
		info.CommunityDocumentation.HasPullRequestTemplate,
		info.CommunityDocumentation.HasLicense,
		info.CommunityDocumentation.HasReadme,
		info.CommunityDocumentation.HasContentReportsEnabled,
		info.CommunityDocumentation.HasWiki)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
