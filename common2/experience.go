package common2

import (
	"context"
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/stdlib/core"
)

// Experience1 - experience functions struct, a nod to Linus Torvalds and plain C
type Experience1 struct {
	AddIngressInference func(h core.ErrorHandler, origin core.Origin, entry inference1.Entry) (int, *core.Status)
	AddEgressInference  func(h core.ErrorHandler, origin core.Origin, entry inference1.Entry) (int, *core.Status)

	GetRateLimitingAction func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status)
	AddRateLimitingAction func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status

	GetRoutingAction func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status)
	AddRoutingAction func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status

	AddRedirectAction func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status
}

var Exp = func() *Experience1 {
	return &Experience1{
		AddIngressInference: func(h core.ErrorHandler, origin core.Origin, e inference1.Entry) (int, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), addInferenceDuration)
			defer cancel()
			id, status := inference1.AddIngress(ctx, nil, e)
			if !status.OK() {
				h.Handle(status)
			}
			return id, status
		},
		AddEgressInference: func(h core.ErrorHandler, origin core.Origin, e inference1.Entry) (int, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), addInferenceDuration)
			defer cancel()
			id, status := inference1.AddEgress(ctx, nil, e)
			if !status.OK() {
				h.Handle(status)
			}
			return id, status
		},
		GetRateLimitingAction: func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
			defer cancel()
			e, status := action1.GetIngressActiveRateLimiting(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		AddRateLimitingAction: func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddIngressRateLimiting(ctx, origin, action)
			if !status.OK() {
				h.Handle(status)
			}
			return status
		},
		GetRoutingAction: func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
			defer cancel()
			e, status := action1.GetIngressActiveRouting(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		AddRoutingAction: func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddIngressRouting(ctx, origin, action)
			if !status.OK() {
				h.Handle(status)
			}
			return status
		},
		AddRedirectAction: func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddIngressRedirect(ctx, origin, action)
			if !status.OK() {
				h.Handle(status)
			}
			return status
		},
		/*
			ResetRoutingAction: func(h core.ErrorHandler, origin core.Origin, agentId string) *core.Status {
				ctx, cancel := context.WithTimeout(context.Background(), resetDuration)
				defer cancel()
				status := action1.EgressResetRouting(ctx, origin, agentId)
				if !status.OK() {
					h.Handle(status)
				}
				return status
			},

		*/
	}
}()
