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

type guidance struct {
	percentile    func(origin core.Origin, h core.ErrorHandler) percentile1.Entry
	shouldProcess func(origin core.Origin, h core.ErrorHandler) bool
}

func newGuidance() *guidance {
	return &guidance{
		percentile: func(origin core.Origin, h core.ErrorHandler) percentile1.Entry {
			ctx, cancel := context.WithTimeout(context.Background(), percentileDuration)
			defer cancel()
			entry, status := percentile1.Get(ctx, origin)
			h.Handle(status, "")
			return entry
		},
		shouldProcess: func(origin core.Origin, h core.ErrorHandler) bool {
			return true
		},
	}
}

/*
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


*/
