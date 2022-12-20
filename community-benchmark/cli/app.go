package cli

import (
	"community-benchmark/cli/bench"
	"github.com/urfave/cli"
)

func CreateApp() *cli.App {
	app := cli.NewApp()
	app.Name = "GitHub community benchmark"
	app.Usage = "Benchmarks a GitHub repository based on how the community interacts with the repository"

	app.Commands = []cli.Command{
		{
			Name:      "bench",
			Action:    bench.Bench,
			Usage:     "Extracts community data about a list of GitHub repositories",
			ArgsUsage: "<repository link> [<repository link> ...]",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "mysql-export",
					Usage: "If the flag si provided, the collected data will be stored in a MySQL database",
				},
			},
		},
	}

	return app
}
