package performance

const MAX_RETRIES = 10
const DOCKER_PORT = ":25565"

type PerformanceBenchmark struct {
	TimeToBootMs   int64
	TimeToRunMs    int64
	StaticMetrics  StaticMetrics
	RuntimeMetrics RuntimeMetrics
}

type StaticMetrics struct {
	ImageSizeBytes int64 `json:"imageSizeBytes"`
}

type RuntimeMetrics struct {
	TotalMemoryUsageMiB float64 `json:"totalMemoryUsageMiB"`
}
