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
	//percentile func(h core.ErrorHandler, origin core.Origin, curr *resiliency1.Percentile) (*resiliency1.Percentile, *core.Status)
	//redirect   func(h core.ErrorHandler, origin core.Origin) (*resiliency1.Redirect, *core.Status)
	//failover   func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Failover, *core.Status)

	entryCDC    func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCEntry, *core.Status)
	redirectCDC func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCRedirect, *core.Status)
	failoverCDC func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCFailover, *core.Status)

	//ingressAssignment func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Assignment, *core.Status)
	//egressAssignment  func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Assignment, *core.Status)
}

var guide = func() *guidance {
	return &guidance{
		/*
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

		*/
		entryCDC: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCEntry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetEntryCDC(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		redirectCDC: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCRedirect, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetRedirectPlanCDC(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		failoverCDC: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCFailover, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetFailoverPlanCDC(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
	}
}()
