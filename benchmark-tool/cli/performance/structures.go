package performance

const MAX_RETRIES = 10
const DOCKER_PORT = ":25565"

type PerformanceBenchmark struct {
	TimeToBootMs   int64
	TimeToRunMs    int64
	MemoryUsageMiB int
	StaticMetrics  StaticMetrics
}

type StaticMetrics struct {
	ImageSize int64 `json:"imageSize"`
}
