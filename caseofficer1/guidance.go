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
	assignments          func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.HostEntry, resiliency1.LastCDCId, *core.Status)
	newAssignments       func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.HostEntry, *core.Status)
	updatedRedirectPlans func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.RedirectPlan, *core.Status)
	updatedFailoverPlans func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.FailoverPlan, *core.Status)
}

var guide = func() *guidance {
	return &guidance{
		assignments: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.HostEntry, resiliency1.LastCDCId, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, last, status := resiliency1.GetHostEntries(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, last, status
		},
		newAssignments: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.HostEntry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetNewHostEntries(ctx, origin, lastId)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		updatedRedirectPlans: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.RedirectPlan, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetUpdatedRedirectPlans(ctx, origin, lastId)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		updatedFailoverPlans: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.FailoverPlan, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetUpdatedFailoverPlans(ctx, origin, lastId)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
	}
}()
