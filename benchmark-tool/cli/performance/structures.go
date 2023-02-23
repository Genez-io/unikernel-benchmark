package performance

import "time"

type PerformanceBenchmark struct {
	TimeToBoot     time.Duration
	TimeToRun      time.Duration
	MemoryUsageMiB int
}
