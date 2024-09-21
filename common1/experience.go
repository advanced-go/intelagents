package common1

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

// Experience - experience interface, with a nod to Linus Torvalds and plain C
type Experience struct {
	ActiveRateLimitingAction func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status)
	AddRateLimitingAction    func(h core.ErrorHandler, origin core.Origin, action *action1.RateLimiting) *core.Status

	ActiveRoutingAction func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status)
	AddRoutingAction    func(h core.ErrorHandler, origin core.Origin, action *action1.Routing) *core.Status

	AddRedirectAction func(h core.ErrorHandler, origin core.Origin, action *action1.Redirect) *core.Status

	AddInference func(h core.ErrorHandler, origin core.Origin, entry *inference1.Entry) (int, *core.Status)
}

var IngressExperience = func() *Experience {
	return &Experience{
		ActiveRateLimitingAction: func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status) {
			return activeRateLimiting(h, origin, true)
		},
		AddRateLimitingAction: func(h core.ErrorHandler, origin core.Origin, action *action1.RateLimiting) *core.Status {
			return addRateLimiting(h, origin, action, true)
		},
		ActiveRoutingAction: func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status) {
			return activeRouting(h, origin, true)
		},
		AddRoutingAction: func(h core.ErrorHandler, origin core.Origin, action *action1.Routing) *core.Status {
			return addRouting(h, origin, action, true)
		},
		AddRedirectAction: func(h core.ErrorHandler, origin core.Origin, action *action1.Redirect) *core.Status {
			return addRedirect(h, origin, action, true)
		},
		AddInference: func(h core.ErrorHandler, origin core.Origin, e *inference1.Entry) (int, *core.Status) {
			return addInference(h, origin, e, true)
		},
	}
}()

var EgressExperience = func() *Experience {
	return &Experience{
		ActiveRateLimitingAction: func(h core.ErrorHandler, origin core.Origin) (action1.RateLimiting, *core.Status) {
			return activeRateLimiting(h, origin, false)
		},
		AddRateLimitingAction: func(h core.ErrorHandler, origin core.Origin, action *action1.RateLimiting) *core.Status {
			return addRateLimiting(h, origin, action, false)
		},
		ActiveRoutingAction: func(h core.ErrorHandler, origin core.Origin) (action1.Routing, *core.Status) {
			return activeRouting(h, origin, false)
		},
		AddRoutingAction: func(h core.ErrorHandler, origin core.Origin, action *action1.Routing) *core.Status {
			return addRouting(h, origin, action, false)
		},
		AddRedirectAction: func(h core.ErrorHandler, origin core.Origin, action *action1.Redirect) *core.Status {
			return addRedirect(h, origin, action, false)
		},
		AddInference: func(h core.ErrorHandler, origin core.Origin, e *inference1.Entry) (int, *core.Status) {
			return addInference(h, origin, e, false)
		},
	}
}()

func activeRateLimiting(h core.ErrorHandler, origin core.Origin, ingress bool) (action1.RateLimiting, *core.Status) {
	ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
	defer cancel()
	e, status := action1.ActiveRateLimiting(ctx, origin, ingress)
	if !status.OK() && !status.NotFound() {
		h.Handle(status)
	}
	return e, status
}

func addRateLimiting(h core.ErrorHandler, origin core.Origin, action *action1.RateLimiting, ingress bool) *core.Status {
	ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
	defer cancel()
	status := action1.AddRateLimiting(ctx, origin, action, ingress)
	if !status.OK() {
		h.Handle(status)
	}
	return status
}

func activeRouting(h core.ErrorHandler, origin core.Origin, ingress bool) (action1.Routing, *core.Status) {
	ctx, cancel := context.WithTimeout(context.Background(), getActionDuration)
	defer cancel()
	e, status := action1.ActiveRouting(ctx, origin, ingress)
	if !status.OK() && !status.NotFound() {
		h.Handle(status)
	}
	return e, status
}

func addRouting(h core.ErrorHandler, origin core.Origin, action *action1.Routing, ingress bool) *core.Status {
	ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
	defer cancel()
	status := action1.AddRouting(ctx, origin, action, ingress)
	if !status.OK() {
		h.Handle(status)
	}
	return status
}

func addRedirect(h core.ErrorHandler, origin core.Origin, action *action1.Redirect, ingress bool) *core.Status {
	ctx, cancel := context.WithTimeout(context.Background(), addActionDuration)
	defer cancel()
	status := action1.AddRedirect(ctx, origin, action, ingress)
	if !status.OK() {
		h.Handle(status)
	}
	return status
}

func addInference(h core.ErrorHandler, origin core.Origin, e *inference1.Entry, ingress bool) (int, *core.Status) {
	ctx, cancel := context.WithTimeout(context.Background(), addInferenceDuration)
	defer cancel()
	id, status := inference1.Add(ctx, origin, e, ingress)
	if !status.OK() {
		h.Handle(status)
	}
	return id, status
}
