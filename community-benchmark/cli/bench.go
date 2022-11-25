package cli

import (
	"context"
	"fmt"
	"github.com/google/go-github/v48/github"
	"github.com/urfave/cli"
	"log"
	"regexp"
	"strings"
)

func bench(c *cli.Context) error {
	for _, arg := range c.Args() {
		client := github.NewClient(nil)

		owner, repo, err := extractRepo(arg)
		if err != nil {
			log.Print(err)
			continue
		}

		list, resp, err := client.PullRequests.List(context.Background(), owner, repo, &github.PullRequestListOptions{ListOptions: github.ListOptions{PerPage: 1000}})
		if err != nil {
			if resp.StatusCode == 404 {
				log.Printf("Repository %s not found. Skipping...", arg)
				continue
			}
			return err
		}

		fmt.Println(len(list))
		for _, pr := range list {
			fmt.Println(*pr.Title)
		}
	}

	return nil
}

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
