package community

import (
	"fmt"
	"log"

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

	for _, arg := range c.Args() {
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
