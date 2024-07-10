package ingress1

import (
	"context"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	percentileDuration = time.Second * 2
	scheduleDuration   = time.Second * 2
)

// getPercentile - resource GET
func getPercentile(origin core.Origin, h core.ErrorHandler) percentile1.Entry {
	ctx, cancel := context.WithTimeout(context.Background(), percentileDuration)
	defer cancel()
	entry, status := percentile1.Get(ctx, origin)
	h.Handle(status, "")
	return entry
}

func processing(origin core.Origin, h core.ErrorHandler) bool {
	return true
}
