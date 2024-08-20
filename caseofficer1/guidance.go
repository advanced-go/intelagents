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
	// TODO : need to distinguish between ingress and egress for assignments??
	assignments          func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Entry, *core.Status)
	newAssignments       func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Entry, *core.Status)
	updatedRedirectPlans func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.RedirectPlan, *core.Status)
	updatedFailoverPlans func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.FailoverPlan, *core.Status)

	//ingressAssignment func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Assignment, *core.Status)
	//egressAssignment  func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Assignment, *core.Status)
}

var guide = func() *guidance {
	return &guidance{
		assignments: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetAssignments(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		newAssignments: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCEntry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetEntryCDC(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		updatedRedirectPlans: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCRedirect, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetRedirectPlanCDC(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		updatedFailoverPlans: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.CDCFailover, *core.Status) {
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
