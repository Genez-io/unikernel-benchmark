package performance

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Benchmark(c *cli.Context) error {
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	buildOptions := types.ImageBuildOptions{
		Remove: true,
	}
	if c.Bool("clear-cache") {
		buildOptions.ForceRemove = true
		buildOptions.NoCache = true
	}

	// buildContext, err := archive.TarWithOptions(".", &archive.TarOptions{IncludeFiles: []string{"unikernels", "benchmark-executable", "benchmark-framework"}})
	// if err != nil {
	// 	return err
	// }

	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	for repo, element := range SupportedUnikernels {
		logrus.Infof("Running benchmark for %s", repo)

		benchmark, err := element.(func(*client.Client, types.ImageBuildOptions) (*PerformanceBenchmark, error))(dockerClient, buildOptions)
		if err != nil {
			return err
		}

		if benchmark != nil {
			logrus.Infof("Time to run: %dms", benchmark.TimeToRunMs)
			logrus.Infof("Time to boot: %dms", benchmark.TimeToBootMs)
			logrus.Infof("Image size: %.2fMiB", float64(benchmark.StaticMetrics.ImageSizeBytes)/(1024*1024))
			logrus.Infof("Total memory usage: %.2fMiB", benchmark.RuntimeMetrics.TotalMemoryUsageMiB)
		}
	}

	return nil
}
