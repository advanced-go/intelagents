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
	percentileSLO  func(h core.ErrorHandler, origin core.Origin) (resiliency1.PercentileSLO, *core.Status)
	redirectPlan   func(h core.ErrorHandler, origin core.Origin) (resiliency1.RedirectPlan, *core.Status)
	updateRedirect func(h core.ErrorHandler, origin core.Origin, status string) *core.Status

	redirectState   func(h core.ErrorHandler, origin core.Origin) (*resiliency1.IngressRedirectState, *core.Status)
	resiliencyState func(h core.ErrorHandler, origin core.Origin) (*resiliency1.IngressResiliencyState, *core.Status)
}

var localGuidance = func() *guidance {
	return &guidance{
		percentileSLO: func(h core.ErrorHandler, origin core.Origin) (resiliency1.PercentileSLO, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetPercentileSLO(ctx, origin)
			if !status.OK() {
				h.Handle(status, "")
			}
			return e, status
		},
		redirectPlan: func(h core.ErrorHandler, origin core.Origin) (resiliency1.RedirectPlan, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetRedirectPlan(ctx, origin)
			if status.OK() {
				return e, status
			}
			if !status.NotFound() {
				h.Handle(status, "")
			}
			return resiliency1.RedirectPlan{}, status
		},
		updateRedirect: func(h core.ErrorHandler, origin core.Origin, status string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addDuration)
			defer cancel()
			status1 := resiliency1.UpdateRedirectPlan(ctx, origin, status)
			if !status1.OK() {
				h.Handle(status1, "")
			}
			return status1
		},
		redirectState: func(h core.ErrorHandler, origin core.Origin) (*resiliency1.IngressRedirectState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			s, status := resiliency1.GetIngressRedirectState(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return s, status
		},
		resiliencyState: func(h core.ErrorHandler, origin core.Origin) (*resiliency1.IngressResiliencyState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			s, status := resiliency1.GetIngressResiliencyState(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return s, status
		},
	}
}()
