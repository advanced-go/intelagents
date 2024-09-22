package common

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
	resetDuration        = time.Second * 10
)

// Experience - experience functions struct, a nod to Linus Torvalds and plain C
type Experience struct {
	AddIngressInference func(h core.ErrorHandler, origin core.Origin, entry inference1.Entry) (int, *core.Status)

	AddRateLimitingAction func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status
	AddRoutingAction      func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status
	AddRedirectAction     func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status

	ResetRoutingAction func(h core.ErrorHandler, origin core.Origin, agentId string) *core.Status
}

var Exp = func() *Experience {
	return &Experience{
		AddIngressInference: func(h core.ErrorHandler, origin core.Origin, e inference1.Entry) (int, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), addInferenceDuration)
			defer cancel()
			id, status := inference1.Add(ctx, origin, &e, true)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return id, status
		},
		AddRateLimitingAction: func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddRateLimiting(ctx, origin, &action, true)
			if !status.OK() {
				h.Handle(status)
			}
			return status
		},
		AddRoutingAction: func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddRouting(ctx, origin, &action, true)
			if !status.OK() {
				h.Handle(status)
			}
			return status
		},
		AddRedirectAction: func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddRedirect(ctx, origin, &action, true)
			if !status.OK() {
				h.Handle(status)
			}
			return status
		},
		/*
			ResetRoutingAction: func(h core.ErrorHandler, origin core.Origin, agentId string) *core.Status {
				ctx, cancel := context.WithTimeout(context.Background(), resetDuration)
				defer cancel()
				status := action1.ResetRouting(ctx, origin, agentId)
				if !status.OK() {
					h.Handle(status)
				}
				return status
			},

		*/
	}
}()
