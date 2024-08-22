package common

import (
	"context"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	addActionDuration = time.Second * 2
)

// Guidance - guidance functions
type Guidance struct {
	AddRateLimitingAction func(h core.ErrorHandler, origin core.Origin, action *resiliency1.RateLimitingAction) *core.Status
	AddRoutingAction      func(h core.ErrorHandler, origin core.Origin, action *resiliency1.RoutingAction) *core.Status
	AddRedirectAction     func(h core.ErrorHandler, origin core.Origin, action *resiliency1.RedirectAction) *core.Status
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
	}
}()
