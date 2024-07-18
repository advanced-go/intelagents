package ingress1

import (
	"context"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

// A nod to Linus Torvalds and plain C
type guidance struct {
	percentile func(duration time.Duration, curr percentile1.Entry, origin core.Origin) (percentile1.Entry, *core.Status)
}

func newGuidance() *guidance {
	return &guidance{
		percentile: func(duration time.Duration, curr percentile1.Entry, origin core.Origin) (percentile1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()
			e, status := percentile1.Get(ctx, origin)
			if status.OK() {
				return e, status
			}
			return curr, status
		},
	}
}
