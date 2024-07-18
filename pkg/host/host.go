package host

import (
	"fmt"
	"github.com/elastic/gosigar"
	"math"
	"runtime"
)

func GetCpu() int {
	return runtime.GOMAXPROCS(0)
}

// GetTotalMemory gets total system memory in bytes
func GetMemory() (int64, error) {
	mem := gosigar.Mem{}
	if err := mem.Get(); err != nil {
		return 0, err
	}
	if mem.Total > math.MaxInt64 {
		return 0, fmt.Errorf("inferred memory size exceeds maximum supported memory size")
	}
	return int64(mem.Total), nil
}
