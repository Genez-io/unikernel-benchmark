package performance

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/urfave/cli"
	"time"
)

func OSvDriver(c *cli.Context) (*PerformanceBenchmark, error) {
	dockerClient, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	buildOptions := types.ImageBuildOptions{
		Tags:       []string{"osv23"},
		Dockerfile: "osv.Dockerfile",
	}

	if c.Bool("clear-cache") {
		buildOptions.ForceRemove = true
		buildOptions.NoCache = true
	}

	tar, err := archive.TarWithOptions("unikernels/", &archive.TarOptions{})
	if err != nil {
		return nil, err
	}

	builtImage, err := dockerClient.ImageBuild(context.Background(), tar, buildOptions)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(builtImage.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	start := time.Now()

	end := time.Now()

	return &PerformanceBenchmark{
		TimeToBootMs:   0,
		TimeToRunMs:    end.Sub(start).Milliseconds(),
		MemoryUsageMiB: 0,
	}, nil
}

func UnikraftDriver(c *cli.Context) (*PerformanceBenchmark, error) {
	return nil, nil
}
