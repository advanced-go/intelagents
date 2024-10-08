package module

import (
	"context"
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	//queryInferenceDuration  = time.Second * 2
	addInferenceDuration = time.Second * 2
	//insertActionDuration    = time.Second * 2
	getActionDuration = time.Second * 2
	addActionDuration = time.Second * 2
)

// A nod to Linus Torvalds and plain C
type experience struct {
	addInference func(h core.ErrorHandler, origin core.Origin, entry inference1.Entry) *core.Status

	getRateLimitingAction func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status)
	getRoutingAction      func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status)

	addRateLimitingAction func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status
	addRoutingAction      func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status
	addRedirectAction     func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status
}

var exp = func() *experience {
	return &experience{
		addInference: func(h core.ErrorHandler, origin core.Origin, e inference1.Entry) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addInferenceDuration)
			defer cancel()
			status := inference1.IngressInsert(ctx, nil, e)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return status
		},
		getRateLimitingAction: func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
			defer cancel()
			action, status := action1.GetRateLimiting(ctx, origin)
			if !status.OK() {
				h.Handle(status, "")
			}
			return action, status
		},
		getRoutingAction: func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
			defer cancel()
			action, status := action1.GetRouting(ctx, origin)
			if !status.OK() {
				h.Handle(status, "")
			}
			return action, status
		},
		addRateLimitingAction: func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddRateLimiting(ctx, origin, action)
			if !status.OK() {
				h.Handle(status, "")
			}
			return status
		},
		addRoutingAction: func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddRouting(ctx, origin, action)
			if !status.OK() {
				h.Handle(status, "")
			}
			return status
		},
		addRedirectAction: func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
			defer cancel()
			status := action1.AddRedirect(ctx, origin, action)
			if !status.OK() {
				h.Handle(status, "")
			}
			return status
		},
		//processControllerAction: controllerAction,
		//reviseTicker:  controllerReviseTicker,
	}
}()
