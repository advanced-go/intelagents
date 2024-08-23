package common

import (
	"context"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	addActionDuration = time.Second * 2
	getDuration       = time.Second * 2
	addDuration       = time.Second * 2
	deleteDuration    = time.Second * 2
)

// Guidance - guidance functions struct, a nod to Linus Torvalds and plain C
type Guidance struct {
	PercentileSLO      func(h core.ErrorHandler, origin core.Origin) (resiliency1.PercentileSLO, *core.Status)
	UpdateRedirect     func(h core.ErrorHandler, origin core.Origin, status string) *core.Status
	DeleteFailoverPlan func(h core.ErrorHandler, origin core.Origin) *core.Status

	RedirectState   func(h core.ErrorHandler, origin core.Origin) (resiliency1.IngressRedirectState, *core.Status)
	ResiliencyState func(h core.ErrorHandler, origin core.Origin) (resiliency1.IngressResiliencyState, *core.Status)
	EgressState     func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.EgressState, *core.Status)

	Assignments          func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.HostEntry, resiliency1.LastCDCId, *core.Status)
	NewAssignments       func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.HostEntry, *core.Status)
	UpdatedRedirectPlans func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.RedirectPlan, *core.Status)
	UpdatedFailoverPlans func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.FailoverPlan, *core.Status)
}

var Guide = func() *Guidance {
	return &Guidance{
		PercentileSLO: func(h core.ErrorHandler, origin core.Origin) (resiliency1.PercentileSLO, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetPercentileSLO(ctx, origin)
			if !status.OK() {
				h.Handle(status, "")
			}
			return e, status
		},
		UpdateRedirect: func(h core.ErrorHandler, origin core.Origin, status string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addDuration)
			defer cancel()
			status1 := resiliency1.UpdateRedirectPlan(ctx, origin, status)
			if !status1.OK() {
				h.Handle(status1, "")
			}
			return status1
		},
		DeleteFailoverPlan: func(h core.ErrorHandler, origin core.Origin) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), deleteDuration)
			defer cancel()
			status := resiliency1.DeleteFailoverPlan(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return status
		},
		RedirectState: func(h core.ErrorHandler, origin core.Origin) (resiliency1.IngressRedirectState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			s, status := resiliency1.GetIngressRedirectState(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return s, status
		},
		ResiliencyState: func(h core.ErrorHandler, origin core.Origin) (resiliency1.IngressResiliencyState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			s, status := resiliency1.GetIngressResiliencyState(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return s, status
		},
		EgressState: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.EgressState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			s, status := resiliency1.GetEgressState(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return s, status
		},
		Assignments: func(h core.ErrorHandler, origin core.Origin) ([]resiliency1.HostEntry, resiliency1.LastCDCId, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, last, status := resiliency1.GetHostEntries(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, last, status
		},
		NewAssignments: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.HostEntry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetNewHostEntries(ctx, origin, lastId)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		UpdatedRedirectPlans: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.RedirectPlan, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := resiliency1.GetUpdatedRedirectPlans(ctx, origin, lastId)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		UpdatedFailoverPlans: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]resiliency1.FailoverPlan, *core.Status) {
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
