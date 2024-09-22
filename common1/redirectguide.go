package common1

import (
	"context"
	"github.com/advanced-go/guidance/redirect1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	getDuration = time.Second * 2
	addDuration = time.Second * 2
)

// RedirectGuidance - redirect guidance interface, with a nod to Linus Torvalds and plain C
type RedirectGuidance struct {
	Ingress   func(h core.ErrorHandler, origin core.Origin) (redirect1.IngressEntry, *core.Status)
	Egress    func(h core.ErrorHandler, origin core.Origin) (redirect1.EgressEntry, *core.Status)
	AllEgress func(h core.ErrorHandler, origin core.Origin) ([]redirect1.EgressEntry, *core.Status)

	AddIngressStatus func(h core.ErrorHandler, origin core.Origin, status, comment string) *core.Status
	AddEgressStatus  func(h core.ErrorHandler, origin core.Origin, status, comment string) *core.Status
}

var RedirectGuide = func() *RedirectGuidance {
	return &RedirectGuidance{
		Ingress: func(h core.ErrorHandler, origin core.Origin) (redirect1.IngressEntry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := redirect1.Ingress(ctx, origin)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		Egress: func(h core.ErrorHandler, origin core.Origin) (redirect1.EgressEntry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), getDuration)
			defer cancel()
			e, status1 := redirect1.Egress(ctx, origin)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		AllEgress: func(h core.ErrorHandler, origin core.Origin) ([]redirect1.EgressEntry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, status1 := redirect1.AllEgress(ctx, origin)
			if !status1.OK() {
				h.Handle(status1)
			}
			return e, status1
		},
		AddIngressStatus: func(h core.ErrorHandler, origin core.Origin, status, comment string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addDuration)
			defer cancel()
			status1 := redirect1.AddStatus(ctx, origin, status, comment, true)
			if !status1.OK() {
				h.Handle(status1)
			}
			return status1
		},
		AddEgressStatus: func(h core.ErrorHandler, origin core.Origin, status, comment string) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addDuration)
			defer cancel()
			status1 := redirect1.AddStatus(ctx, origin, status, comment, false)
			if !status1.OK() {
				h.Handle(status1)
			}
			return status1
		},
	}
}()
