package community

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Benchmark(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("no GitHub repositories provided")
	}

	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	for _, arg := range c.Args() {
		logrus.Info("Collecting community data for ", bold(green(arg)))
		owner, repo, err := extractRepo(arg)
		if err != nil {
			log.Print(err)
			continue
		}

		info, err := collectRepositoryInfo(owner, repo)
		if err != nil {
			log.Print(err)
			continue
		}

		err = saveOutput(c, info)
		if err != nil {
			log.Print(err)
			continue
		}
	}

	return nil
}
