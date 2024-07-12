package guidance1

import (
	"time"
)

const (
	PercentilePollingDuration = time.Hour * 12
)

func IsScheduled() bool {
	return true
}
