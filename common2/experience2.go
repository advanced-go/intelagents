package common2

import (
	"context"
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	addInferenceDuration = time.Second * 2
	addActionDuration    = time.Second * 2
	getActionDuration    = time.Second * 2
	resetDuration        = time.Second * 10
)

// Experience - experience functions struct, a nod to Linus Torvalds and plain C
type Experience struct {
	GetLastRateLimitingAction func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status)
	AddRateLimitingAction     func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status
	GetLastRoutingAction      func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status)
	AddRoutingAction          func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status
	AddRedirectAction         func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status

	AddInference func(h core.ErrorHandler, origin core.Origin, entry inference1.Entry) (int, *core.Status)
}

var IngressExperience = func() *Experience {
	return &Experience{
		GetLastRateLimitingAction: func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status) {
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
		GetLastRoutingAction: func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status) {
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
		AddInference: func(h core.ErrorHandler, origin core.Origin, e inference1.Entry) (int, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), addInferenceDuration)
			defer cancel()
			id, status := inference1.AddIngress(ctx, nil, e)
			if !status.OK() {
				h.Handle(status)
			}
			return id, status
		},
	}
}()

var EgressExperience = func() *Experience {
	return &Experience{
		GetLastRateLimitingAction: func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
			defer cancel()
			e, status := action1.GetEgressActiveRateLimiting(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		AddRateLimitingAction: func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddEgressRateLimiting(ctx, origin, action)
			if !status.OK() {
				h.Handle(status)
			}
			return status
		},
		GetLastRoutingAction: func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
			defer cancel()
			e, status := action1.GetEgressActiveRouting(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		AddRoutingAction: func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddEgressRouting(ctx, origin, action)
			if !status.OK() {
				h.Handle(status)
			}
			return status
		},
		AddRedirectAction: func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddEgressRedirect(ctx, origin, action)
			if !status.OK() {
				h.Handle(status)
			}
			return status
		},
		AddInference: func(h core.ErrorHandler, origin core.Origin, e inference1.Entry) (int, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), addInferenceDuration)
			defer cancel()
			id, status := inference1.AddEgress(ctx, nil, e)
			if !status.OK() {
				h.Handle(status)
			}
			return id, status
		},
	}
}()
