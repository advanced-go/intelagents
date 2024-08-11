package ingress1

import (
	"context"
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	//queryInferenceDuration  = time.Second * 2
	insertInferenceDuration = time.Second * 2
	insertActionDuration    = time.Second * 2
)

// A nod to Linus Torvalds and plain C
type experience struct {
	addInference     func(h core.ErrorHandler, origin core.Origin, entry inference1.Entry) *core.Status
	processInference func(c *controller, entry []timeseries1.Entry, percentile percentile1.Entry) (inference1.Entry, *core.Status)

	addRateLimitingAction   func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status
	addRoutingAction        func(h core.ErrorHandler, origin core.Origin, action action1.Routing) *core.Status
	addRedirectAction       func(h core.ErrorHandler, origin core.Origin, action action1.Redirect) *core.Status
	processControllerAction func(c *controller, entry inference1.Entry) (action1.RateLimiting, *core.Status)
	processRoutingAction    func(c *controller, entry inference1.Entry) (action1.Routing, *core.Status)
	//processRedirectAction   func(c *controller, entry inference1.Entry) (action1.Redirect, *core.Status)

	//reviseTicker func(c *controller)
}

var exp = func() *experience {
	return &experience{
		addInference: func(h core.ErrorHandler, origin core.Origin, e inference1.Entry) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), insertInferenceDuration)
			defer cancel()
			status := inference1.IngressInsert(ctx, nil, e)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return status
		},
		processInference: controllerInference,
		addRateLimitingAction: func(h core.ErrorHandler, origin core.Origin, action action1.RateLimiting) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), insertActionDuration)
			defer cancel()
			status := action1.AddRateLimiting(ctx, origin, action)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return status
		},
		//processControllerAction: controllerAction,
		//reviseTicker:  controllerReviseTicker,
	}
}()
