package common

import (
	"context"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	addInferenceDuration = time.Second * 2
)

// A nod to Linus Torvalds and plain C

// Experience - experience functions
type Experience struct {
	AddInference func(h core.ErrorHandler, origin core.Origin, entry inference1.Entry) *core.Status

	/*
		GetRateLimitingAction func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status)
		GetRoutingAction      func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status)

		AddRateLimitingAction func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status
		AddRoutingAction      func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status
		AddRedirectAction     func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status

	*/
}

var Exp = func() *Experience {
	return &Experience{
		AddInference: func(h core.ErrorHandler, origin core.Origin, e inference1.Entry) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addInferenceDuration)
			defer cancel()
			status := inference1.IngressInsert(ctx, nil, e)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return status
		},
		/*
			GetRateLimitingAction: func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status) {
				ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
				defer cancel()
				action, status := action1.GetRateLimiting(ctx, origin)
				if !status.OK() {
					h.Handle(status, "")
				}
				return action, status
			},
			GetRoutingAction: func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status) {
				ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
				defer cancel()
				action, status := action1.GetRouting(ctx, origin)
				if !status.OK() {
					h.Handle(status, "")
				}
				return action, status
			},
			AddRateLimitingAction: func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status {
				ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
				defer cancel()
				status := action1.AddRateLimiting(ctx, origin, action)
				if !status.OK() {
					h.Handle(status, "")
				}
				return status
			},
			AddRoutingAction: func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status {
				ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
				defer cancel()
				status := action1.AddRouting(ctx, origin, action)
				if !status.OK() {
					h.Handle(status, "")
				}
				return status
			},
			AddRedirectAction: func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status {
				ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
				defer cancel()
				status := action1.AddRedirect(ctx, origin, action)
				if !status.OK() {
					h.Handle(status, "")
				}
				return status
			},

		*/
	}
}()
