package ingress1

import (
	"context"
	"github.com/advanced-go/guidance/controller1"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/guidance/schedule1"
	"github.com/advanced-go/guidance/update1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	controllerDuration = time.Second * 2
	versionDuration    = time.Second * 2
	updateDuration     = time.Second * 2
)

// A nod to Linus Torvalds and plain C
type guidance struct {
	isScheduled       func(origin core.Origin) bool
	percentile        func(duration time.Duration, curr percentile1.Entry, origin core.Origin) (percentile1.Entry, *core.Status)
	controllers       func(origin core.Origin) (controller1.Ingress, *core.Status)
	controllerVersion func(origin core.Origin) (controller1.Entry, *core.Status)
	updateRedirect    func(origin core.Origin, status string) *core.Status
}

func newGuidance(agent messaging.OpsAgent) *guidance {
	return &guidance{
		isScheduled: func(origin core.Origin) bool {
			return schedule1.IsIngressControllerScheduled(origin)
		},
		percentile: func(duration time.Duration, curr percentile1.Entry, origin core.Origin) (percentile1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()
			e, status := percentile1.Get(ctx, origin)
			if status.OK() {
				return e, status
			}
			agent.Handle(status, "")
			return curr, status
		},
		controllers: func(origin core.Origin) (controller1.Ingress, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), controllerDuration)
			defer cancel()
			e, status := controller1.IngressControllers(ctx, origin)
			if status.OK() {
				return e[0], status
			}
			if !status.NotFound() {
				agent.Handle(status, "")
			}
			return controller1.Ingress{}, status
		},
		controllerVersion: func(origin core.Origin) (controller1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), versionDuration)
			defer cancel()
			e, status := controller1.Version(ctx, origin)
			if status.OK() || status.NotFound() {
				return e, status
			}
			agent.Handle(status, "")
			return controller1.Entry{}, status
		},
		updateRedirect: func(origin core.Origin, status string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), updateDuration)
			defer cancel()
			status1 := update1.IngressRedirect(ctx, origin, status)
			if !status1.OK() {
				agent.Handle(status1, "")
			}
			return status1
		},
	}
}
