package community

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
)

func Benchmark(c *cli.Context) error {
	if c.NArg() == 0 {
		return fmt.Errorf("no GitHub repositories provided")
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
