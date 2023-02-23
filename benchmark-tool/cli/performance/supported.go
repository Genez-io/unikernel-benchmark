package performance

import (
	"benchmark-tool/cli/performance/drivers"
)

var supportedUnikernels = map[string]interface{}{
	"cloudius-systems/osv": drivers.OSvDriver,
	"unikraft/unikraft":    drivers.UnikraftDriver,
}
