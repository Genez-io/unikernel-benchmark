package performance

import (
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/urfave/cli"
)

func Benchmark(c *cli.Context) error {
	buildOptions := types.ImageBuildOptions{}
	if c.Bool("clear-cache") {
		buildOptions.ForceRemove = true
		buildOptions.NoCache = true
	}

	buildContext, err := archive.TarWithOptions(".", &archive.TarOptions{IncludeFiles: []string{"unikernels", "benchmark-executable", "benchmark-framework"}})
	if err != nil {
		return err
	}

	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	for _, element := range SupportedUnikernels {
		benchmark, err := element.(func(*client.Client, io.ReadCloser, types.ImageBuildOptions) (*PerformanceBenchmark, error))(dockerClient, buildContext, buildOptions)
		if err != nil {
			return err
		}

		if benchmark != nil {
			log.Printf("Memory used: %dMiB", benchmark.MemoryUsageMiB)
			log.Printf("Time to run: %dms", benchmark.TimeToRunMs)
			log.Printf("Time to boot: %dms", benchmark.TimeToBootMs)
		}
	}

	return nil
}
