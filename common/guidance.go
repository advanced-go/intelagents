package common

import (
	"context"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	addActionDuration = time.Second * 2
	getDuration       = time.Second * 2
	addDuration       = time.Second * 2
)

// Guidance - guidance functions struct, a nod to Linus Torvalds and plain C
type Guidance struct {
	AddRateLimitingAction func(h core.ErrorHandler, origin core.Origin, action *resiliency1.RateLimitingAction) *core.Status
	AddRoutingAction      func(h core.ErrorHandler, origin core.Origin, action *resiliency1.RoutingAction) *core.Status
	AddRedirectAction     func(h core.ErrorHandler, origin core.Origin, action *resiliency1.RedirectAction) *core.Status

	PercentileSLO func(h core.ErrorHandler, origin core.Origin) (resiliency1.PercentileSLO, *core.Status)

	UpdateRedirect func(h core.ErrorHandler, origin core.Origin, status string) *core.Status

	FailoverPlan func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.FailoverPlan, *core.Status)

	RedirectState   func(h core.ErrorHandler, origin core.Origin) (*resiliency1.IngressRedirectState, *core.Status)
	ResiliencyState func(h core.ErrorHandler, origin core.Origin) (*resiliency1.IngressResiliencyState, *core.Status)
	EgressState     func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.EgressState, *core.Status)
}

var Guide = func() *Guidance {
	return &Guidance{
		AddRateLimitingAction: func(h core.ErrorHandler, origin core.Origin, action *resiliency1.RateLimitingAction) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := resiliency1.AddRateLimitingAction(ctx, origin, action)
			if !status.OK() {
				h.Handle(status, "")
			}
			return status
		},
		AddRoutingAction: func(h core.ErrorHandler, origin core.Origin, action *resiliency1.RoutingAction) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := resiliency1.AddRoutingAction(ctx, origin, action)
			if !status.OK() {
				h.Handle(status, "")
			}
			return status
		},
		AddRedirectAction: func(h core.ErrorHandler, origin core.Origin, action *resiliency1.RedirectAction) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := resiliency1.AddRedirectAction(ctx, origin, action)
			if !status.OK() {
				h.Handle(status, "")
			}
			return status
		},
		PercentileSLO: func(h core.ErrorHandler, origin core.Origin) (resiliency1.PercentileSLO, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetPercentileSLO(ctx, origin)
			if !status.OK() {
				h.Handle(status, "")
			}
			return e, status
		},
		UpdateRedirect: func(h core.ErrorHandler, origin core.Origin, status string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addDuration)
			defer cancel()
			status1 := resiliency1.UpdateRedirectPlan(ctx, origin, status)
			if !status1.OK() {
				h.Handle(status1, "")
			}
			return status1
		},
		FailoverPlan: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.FailoverPlan, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			s, status := resiliency1.GetFailoverPlan(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return s, status
		},
		RedirectState: func(h core.ErrorHandler, origin core.Origin) (*resiliency1.IngressRedirectState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			s, status := resiliency1.GetIngressRedirectState(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return s, status
		},
		ResiliencyState: func(h core.ErrorHandler, origin core.Origin) (*resiliency1.IngressResiliencyState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			s, status := resiliency1.GetIngressResiliencyState(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return s, status
		},
		EgressState: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.EgressState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			s, status := resiliency1.GetEgressState(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return s, status
		},
	}
}()
