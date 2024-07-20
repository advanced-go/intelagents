package ingress1

import (
	"context"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/guidance/schedule1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

// A nod to Linus Torvalds and plain C
type guidance struct {
	isScheduled func(origin core.Origin) bool
	percentile  func(duration time.Duration, curr percentile1.Entry, origin core.Origin) (percentile1.Entry, *core.Status)
}

func newGuidance(agent messaging.OpsAgent) *guidance {
	return &guidance{
		isScheduled: func(origin core.Origin) bool {
			return schedule1.IsIngressControllerScheduled(origin)
		},
		percentile: func(duration time.Duration, curr percentile1.Entry, origin core.Origin) (percentile1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()
			e, status := percentile1.Get(ctx, origin)
			if status.OK() {
				return e, status
			}
			agent.Handle(status, "")
			return curr, status
		},
	}
}
