package common2

import (
	"context"
	"github.com/advanced-go/guidance/host1"
	"github.com/advanced-go/guidance/redirect1"
	"github.com/advanced-go/guidance/state1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	getDuration    = time.Second * 2
	addDuration    = time.Second * 2
	deleteDuration = time.Second * 2
)

// Guidance - guidance functions struct, a nod to Linus Torvalds and plain C
type Guidance struct {
	HostEntries    func(h core.ErrorHandler, origin core.Origin) ([]host1.Entry, *core.Status)
	NewHostEntries func(h core.ErrorHandler, origin core.Origin, lastId int) ([]host1.Entry, *core.Status)

	IngressRedirect          func(h core.ErrorHandler, origin core.Origin) ([]redirect1.RedirectConfig, *core.Status)
	UpdatedIngressRedirect   func(h core.ErrorHandler, origin core.Origin, lastId int) ([]redirect1.RedirectConfig, *core.Status)
	AddIngressRedirectStatus func(h core.ErrorHandler, origin core.Origin, status string) *core.Status

	EgressRedirect        func(h core.ErrorHandler, origin core.Origin) ([]redirect1.RedirectConfig, *core.Status)
	UpdatedEgressRedirect func(h core.ErrorHandler, origin core.Origin, lastId int) ([]redirect1.RedirectConfig, *core.Status)

	IngressRedirectState     func(h core.ErrorHandler, origin core.Origin) (state1.IngressRedirectState, *core.Status)
	IngressRateLimitingState func(h core.ErrorHandler, origin core.Origin) (state1.IngressResiliencyState, *core.Status)

	EgressRedirectState     func(h core.ErrorHandler, origin core.Origin) (state1.IngressRedirectState, *core.Status)
	EgressRateLimitingState func(h core.ErrorHandler, origin core.Origin) ([]state1.EgressState, *core.Status)
}

var Guide = func() *Guidance {
	return &Guidance{
		HostEntries: func(h core.ErrorHandler, origin core.Origin) ([]host1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, _, status := host1.GetHostEntries(ctx, origin)
			if !status.OK() {
				h.Handle(status)
			}
			return e, status
		},
		NewHostEntries: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]host1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status := host1.GetNewHostEntries(ctx, origin, lastId)
			if !status.OK() {
				h.Handle(status)
			}
			return e, status
		},
		IngressRedirect: func(h core.ErrorHandler, origin core.Origin) ([]redirect1.RedirectConfig, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := redirect1.GetIngressRedirect(ctx, origin)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		UpdatedIngressRedirect: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]redirect1.RedirectConfig, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := redirect1.GetUpdatedIngressRedirect(ctx, origin, lastId)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		AddIngressRedirectStatus: func(h core.ErrorHandler, origin core.Origin, status string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addDuration)
			defer cancel()
			status1 := redirect1.AddIngressRedirectStatus(ctx, origin, status)
			if !status1.OK() {
				h.Handle(status1)
			}
			return status1
		},
		EgressRedirect: func(h core.ErrorHandler, origin core.Origin) ([]redirect1.RedirectConfig, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := redirect1.GetEgressRedirect(ctx, origin)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		UpdatedEgressRedirect: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]redirect1.RedirectConfig, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := redirect1.GetUpdatedEgressRedirect(ctx, origin, lastId)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		IngressRedirectState: func(h core.ErrorHandler, origin core.Origin) (state1.IngressRedirectState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := state1.GetIngressRedirectState(ctx, origin)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		IngressRateLimitingState: func(h core.ErrorHandler, origin core.Origin) (state1.IngressResiliencyState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := state1.GetIngressRateLimitingState(ctx, origin)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},

		EgressRedirectState: func(h core.ErrorHandler, origin core.Origin) (state1.IngressRedirectState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := state1.GetEgressRedirectState(ctx, origin)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		EgressRateLimitingState: func(h core.ErrorHandler, origin core.Origin) ([]state1.EgressState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := state1.GetEgressRateLimitingState(ctx, origin)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
	}
}()
