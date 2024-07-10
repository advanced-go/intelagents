package guidance

import (
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	scheduleDuration = time.Second * 2
)

func init() {
}

func IngressProcessing(origin core.Origin, h core.ErrorHandler) bool {
	return true
}

func EgressProcessing(origin core.Origin, h core.ErrorHandler) bool {
	return true
}
