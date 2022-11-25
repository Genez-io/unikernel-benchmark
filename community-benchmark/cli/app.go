package cli

import (
	"github.com/urfave/cli"
)

func CreateApp() *cli.App {
	app := cli.NewApp()
	app.Name = "GitHub community benchmark"
	app.Usage = "Benchmarks a GitHub repository based on how the community interacts with the repository"

	app.Commands = []cli.Command{
		{
			Name:      "bench",
			Action:    bench,
			Usage:     "Extracts community data about a list of GitHub repositories",
			ArgsUsage: "<repository link> [<repository link> ...]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:      "file-output",
					TakesFile: true,
				},
				&cli.StringFlag{
					Name: "output-format",
					// Usage: "",
					Value: "json",
				},
			},
		},
	}

	return app
}
