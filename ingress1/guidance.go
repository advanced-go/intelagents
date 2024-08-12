package ingress1

import (
	"context"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	addDuration = time.Second * 2
	getDuration = time.Second * 2
)

// A nod to Linus Torvalds and plain C
type guidance struct {
	percentile        func(h core.ErrorHandler, origin core.Origin, curr resiliency1.Percentile) (resiliency1.Percentile, *core.Status)
	redirect          func(h core.ErrorHandler, origin core.Origin) (resiliency1.Redirect, *core.Status)
	addRedirectStatus func(h core.ErrorHandler, origin core.Origin, status string) *core.Status
}

func guide() *guidance {
	return &guidance{
		percentile: func(h core.ErrorHandler, origin core.Origin, curr resiliency1.Percentile) (resiliency1.Percentile, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.IngressPercentile(ctx, origin)
			if status.OK() {
				return e, status
			}
			h.Handle(status, "")
			return curr, status
		},
		redirect: func(h core.ErrorHandler, origin core.Origin) (resiliency1.Redirect, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.IngressRedirect(ctx, origin)
			if status.OK() {
				return e, status
			}
			if !status.NotFound() {
				h.Handle(status, "")
			}
			return resiliency1.Redirect{}, status
		},
		addRedirectStatus: func(h core.ErrorHandler, origin core.Origin, status string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addDuration)
			defer cancel()
			status1 := resiliency1.AddRedirectStatus(ctx, origin, status)
			if !status1.OK() {
				h.Handle(status1, "")
			}
			return status1
		},
	}
}
