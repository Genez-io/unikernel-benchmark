package cli

import (
	"benchmark-tool/cli/community"
	"benchmark-tool/cli/performance"
	"github.com/urfave/cli"
)

func CreateApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Unikernel benchmark tool"
	app.Usage = "Benchmarks a list of supported unikernels"

	app.Commands = []cli.Command{
		{
			Name:      "community",
			Action:    community.Benchmark,
			Usage:     "Benchmarks a list of GitHub repositories based on how the community interacts with them",
			ArgsUsage: "<repository link> [<repository link> ...]",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "mysql-export",
					Usage: "If the flag is provided, the collected data will be stored in a MySQL database",
				},
			},
		},
		{
			Name:      "performance",
			Action:    performance.Benchmark,
			Usage:     "Benchmarks a list of supported unikernels",
			ArgsUsage: "<repository link> [<repository link> ...]",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "clear-cache",
					Usage: "Clear the docker build image cache before building the unikernels",
				},
				&cli.BoolFlag{
					Name:  "mysql-export",
					Usage: "Export the collected data will be stored in a MySQL database",
				},
			},
		},
	}

	return app
}
