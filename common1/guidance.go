package common1

import (
	"context"
	"github.com/advanced-go/guidance/host1"
	"github.com/advanced-go/guidance/redirect1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	getDuration   = time.Second * 2
	addDuration   = time.Second * 2
	queryDuration = time.Second * 3
)

// Guidance - guidance interface, with a nod to Linus Torvalds and plain C
type Guidance struct {
	QueryNewRedirect      func(h core.ErrorHandler, origin core.Origin, lastCDCId int) ([]core.Origin, *core.Status)
	QueryInactiveRedirect func(h core.ErrorHandler, origin core.Origin, lastCDCId int) ([]core.Origin, *core.Status)
	//GetRedirect           func(h core.ErrorHandler, origin core.Origin) (redirect1.Entry, *core.Status)
	//GetHostRedirect       func(h core.ErrorHandler, origin core.Origin) ([]redirect1.Entry, *core.Status)
	AddStatus func(h core.ErrorHandler, origin core.Origin, status string) *core.Status

	/*
		IngressRedirect          func(h core.ErrorHandler, origin core.Origin) ([]redirect1.RedirectConfig, *core.Status)
		UpdatedIngressRedirect   func(h core.ErrorHandler, origin core.Origin, lastId int) ([]redirect1.RedirectConfig, *core.Status)
		AddIngressRedirectStatus func(h core.ErrorHandler, origin core.Origin, status string) *core.Status

		EgressRedirect        func(h core.ErrorHandler, origin core.Origin) ([]redirect1.RedirectConfig, *core.Status)
		UpdatedEgressRedirect func(h core.ErrorHandler, origin core.Origin, lastId int) ([]redirect1.RedirectConfig, *core.Status)

		IngressRedirectState     func(h core.ErrorHandler, origin core.Origin) (state1.IngressRedirectState, *core.Status)
		IngressRateLimitingState func(h core.ErrorHandler, origin core.Origin) (state1.IngressResiliencyState, *core.Status)

		EgressRedirectState     func(h core.ErrorHandler, origin core.Origin) (state1.IngressRedirectState, *core.Status)
		EgressRateLimitingState func(h core.ErrorHandler, origin core.Origin) ([]state1.EgressState, *core.Status)
	*/
}

var IngressGuidance = func() *Guidance {
	return &Guidance{
		QueryNewRedirect: func(h core.ErrorHandler, origin core.Origin, lastCIDId int) ([]core.Origin, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, status1 := redirect1.QueryIngressNew(ctx, origin, lastCIDId)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		QueryInactiveRedirect: func(h core.ErrorHandler, origin core.Origin, lastCIDId int) ([]core.Origin, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, status1 := redirect1.QueryIngressInactive(ctx, origin, lastCIDId)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		/*
			GetRedirect: func(h core.ErrorHandler, origin core.Origin) (redirect1.Entry, *core.Status) {
				ctx, cancel := context.WithTimeout(context.Background(), getDuration)
				defer cancel()
				e, status := redirect1.GetIngress(ctx, origin)
				if !status.OK() {
					h.Handle(status)
				}
				return e, status
			},
			GetHostRedirect: func(_ core.ErrorHandler, _ core.Origin) ([]redirect1.Entry, *core.Status) {
				return nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("error: Ingress - GetHostRedirect() is not implemented"))
			},
			AddStatus: func(h core.ErrorHandler, origin core.Origin, status string) *core.Status {
				ctx, cancel := context.WithTimeout(context.Background(), addDuration)
				defer cancel()
				status1 := redirect1.AddIngressStatus(ctx, origin, status)
				if !status1.OK() {
					h.Handle(status1)
				}
				return status1
			},

		*/
	}
}()

var EgressGuidance = func() *Guidance {
	return &Guidance{
		QueryNewRedirect: func(h core.ErrorHandler, origin core.Origin, lastCIDId int) ([]core.Origin, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, status1 := redirect1.QueryEgressNew(ctx, origin, lastCIDId)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		QueryInactiveRedirect: func(h core.ErrorHandler, origin core.Origin, lastCIDId int) ([]core.Origin, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, status1 := redirect1.QueryEgressInactive(ctx, origin, lastCIDId)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		/*
			GetRedirect: func(h core.ErrorHandler, origin core.Origin) (redirect1.Entry, *core.Status) {
				ctx, cancel := context.WithTimeout(context.Background(), getDuration)
				defer cancel()
				e, status := redirect1.GetEgress(ctx, origin)
				if !status.OK() {
					h.Handle(status)
				}
				return e, status
			},

			GetHostRedirect: func(h core.ErrorHandler, origin core.Origin) ([]redirect1.Entry, *core.Status) {
				ctx, cancel := context.WithTimeout(context.Background(), getDuration)
				defer cancel()
				e, status := redirect1.GetHostEgress(ctx, origin)
				if !status.OK() {
					h.Handle(status)
				}
				return e, status
			},

		*/
		AddStatus: func(h core.ErrorHandler, origin core.Origin, status string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addDuration)
			defer cancel()
			status1 := redirect1.AddEgressStatus(ctx, origin, status)
			if !status1.OK() {
				h.Handle(status1)
			}
			return status1
		},
	}
}()

// hostGuidance - host guidance interface
type hostGuidance struct {
	QueryIngressHosts    func(h core.ErrorHandler, origin core.Origin) ([]host1.EntryQuery, host1.LastCDCId, *core.Status)
	QueryNewIngressHosts func(h core.ErrorHandler, origin core.Origin, lastId int) ([]host1.Entry, host1.LastCDCId, *core.Status)

	QueryEgressHosts    func(h core.ErrorHandler, origin core.Origin) ([]host1.EntryQuery, host1.LastCDCId, *core.Status)
	QueryNewEgressHosts func(h core.ErrorHandler, origin core.Origin, lastId int) ([]host1.Entry, host1.LastCDCId, *core.Status)
}

var HostGuidance = func() *hostGuidance {
	return &hostGuidance{
		QueryIngressHosts: func(h core.ErrorHandler, origin core.Origin) ([]host1.EntryQuery, host1.LastCDCId, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, id, status := host1.QueryIngressHosts(ctx, origin)
			if !status.OK() {
				h.Handle(status)
			}
			return e, id, status
		},
		QueryNewIngressHosts: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]host1.Entry, host1.LastCDCId, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, id, status := host1.QueryNewIngressHosts(ctx, origin, lastId)
			if !status.OK() {
				h.Handle(status)
			}
			return e, id, status
		},
		QueryEgressHosts: func(h core.ErrorHandler, origin core.Origin) ([]host1.EntryQuery, host1.LastCDCId, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, id, status := host1.QueryEgressHosts(ctx, origin)
			if !status.OK() {
				h.Handle(status)
			}
			return e, id, status
		},
		QueryNewEgressHosts: func(h core.ErrorHandler, origin core.Origin, lastId int) ([]host1.Entry, host1.LastCDCId, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, id, status := host1.QueryNewEgressHosts(ctx, origin, lastId)
			if !status.OK() {
				h.Handle(status)
			}
			return e, id, status
		},
	}
}()
