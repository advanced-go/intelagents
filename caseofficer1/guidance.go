package caseofficer1

import (
	"context"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	getDuration = time.Second * 2
)

// A nod to Linus Torvalds and plain C
type guidance struct {
	percentile func(h core.ErrorHandler, origin core.Origin, curr *resiliency1.Percentile) (*resiliency1.Percentile, *core.Status)
	redirect   func(h core.ErrorHandler, origin core.Origin) (*resiliency1.Redirect, *core.Status)
	failover   func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Failover, *core.Status)

	ingressCDC func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCEntry, *core.Status)
	egressCDC  func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCEntry, *core.Status)

	//ingressAssignment func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Assignment, *core.Status)
	//egressAssignment  func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Assignment, *core.Status)
}

var guide = func() *guidance {
	return &guidance{
		percentile: func(h core.ErrorHandler, origin core.Origin, curr *resiliency1.Percentile) (*resiliency1.Percentile, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.IngressPercentile(ctx, origin)
			if status.OK() {
				return &e, status
			}
			h.Handle(status, "")
			return curr, status
		},
		redirect: func(h core.ErrorHandler, origin core.Origin) (*resiliency1.Redirect, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.IngressRedirect(ctx, origin)
			if status.OK() || status.NotFound() {
				return &e, status
			}
			h.Handle(status, "")
			return &resiliency1.Redirect{}, status
		},
		failover: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Failover, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := resiliency1.EgressFailover(ctx, origin)
			if !status1.OK() && !status1.NotFound() {
				h.Handle(status1, "")
			}
			return e, status1
		},
		ingressCDC: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCEntry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := resiliency1.IngressCDC(ctx, origin)
			if !status1.OK() && !status1.NotFound() {
				h.Handle(status1, "")
			}
			return e, status1
		},
		egressCDC: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCEntry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := resiliency1.EgressCDC(ctx, origin)
			if !status1.OK() && !status1.NotFound() {
				h.Handle(status1, "")
			}
			return e, status1
		},
	}
}()
