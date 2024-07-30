package ingress1

import (
	"context"
	"github.com/advanced-go/guidance/controller1"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/guidance/update1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	controllerDuration = time.Second * 2
	versionDuration    = time.Second * 2
	updateDuration     = time.Second * 2
	percentileDuration = time.Second * 2
)

// A nod to Linus Torvalds and plain C
type guidance struct {
	percentile     func(h core.ErrorHandler, origin core.Origin, curr percentile1.Entry) (percentile1.Entry, *core.Status)
	controllers    func(h core.ErrorHandler, origin core.Origin) (controller1.Ingress, *core.Status)
	updateRedirect func(h core.ErrorHandler, origin core.Origin, status string) *core.Status
}

var guide = func() *guidance {
	return &guidance{
		percentile: func(h core.ErrorHandler, origin core.Origin, curr percentile1.Entry) (percentile1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), percentileDuration)
			defer cancel()
			e, status := percentile1.Get(ctx, origin)
			if status.OK() {
				return e, status
			}
			h.Handle(status, "")
			return curr, status
		},
		controllers: func(h core.ErrorHandler, origin core.Origin) (controller1.Ingress, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), controllerDuration)
			defer cancel()
			e, status := controller1.IngressControllers(ctx, origin)
			if status.OK() {
				return e[0], status
			}
			if !status.NotFound() {
				h.Handle(status, "")
			}
			return controller1.Ingress{}, status
		},
		updateRedirect: func(h core.ErrorHandler, origin core.Origin, status string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), updateDuration)
			defer cancel()
			status1 := update1.IngressRedirect(ctx, origin, status)
			if !status1.OK() {
				h.Handle(status1, "")
			}
			return status1
		},
	}
}()
