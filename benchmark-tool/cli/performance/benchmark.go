package performance

import (
	"github.com/urfave/cli"
)

func Benchmark(c *cli.Context) error {
	for _, element := range SupportedUnikernels {
		benchmark, err := element.(func(ctx *cli.Context) (*PerformanceBenchmark, error))(c)
		if err != nil {
			return err
		}

		if benchmark != nil {
			println(benchmark.MemoryUsageMiB)
			println(benchmark.TimeToRunMs)
			println(benchmark.TimeToBootMs)
		}
	}

	return nil
}
