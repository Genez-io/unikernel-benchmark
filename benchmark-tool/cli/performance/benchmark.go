package performance

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
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

	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	for repo, element := range SupportedUnikernels {
		for _, vmm := range element.SupportedVMMs {
			logrus.Infof("Running benchmark for %s with %s backend", bold(green(repo)), bold(red(vmm)))

			benchmark, err := BenchmarkUnikernelWithDocker(dockerClient, buildOptions, element.UnikernelName, vmm)
			if err != nil {
				return err
			}

			if benchmark != nil {
				logrus.Infof("Time to run: %s", bold(blue(fmt.Sprintf("%dms", benchmark.TimeToRunMs))))
				logrus.Infof("Time to boot: %s", bold(blue(fmt.Sprintf("%dms", benchmark.TimeToBootMs))))
				logrus.Infof("Image size: %s", bold(blue(fmt.Sprintf("%.2fMiB", float64(benchmark.StaticMetrics.ImageSizeBytes)/(1024*1024)))))
				logrus.Infof("Total memory usage: %s", bold(blue(fmt.Sprintf("%.2fMiB", benchmark.RuntimeMetrics.TotalMemoryUsageMiB))))
			}
		}
	}

	return nil
}
