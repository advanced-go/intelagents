package common1

import (
	"context"
	"github.com/advanced-go/guidance/host1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	queryDuration = time.Second * 3
)

// hostGuidance - host guidance interface, with a nod to Linus Torvalds and plain C
type hostGuidance struct {
	HostQuery            func(h core.ErrorHandler, origin core.Origin, state *host1.CDCState) ([]host1.EntryStatus, []host1.CDCState, *core.Status)
	RedirectStateChanges func(h core.ErrorHandler, origin core.Origin, state *host1.CDCState, ingress bool) ([]core.Origin, *core.Status)
}

var HostGuidance = func() *hostGuidance {
	return &hostGuidance{
		HostQuery: func(h core.ErrorHandler, origin core.Origin, state *host1.CDCState) ([]host1.EntryStatus, []host1.CDCState, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, cdc, status := host1.HostQuery(ctx, origin, state)
			if !status.OK() {
				h.Handle(status)
			}
			return e, cdc, status
		},
		RedirectStateChanges: func(h core.ErrorHandler, origin core.Origin, state *host1.CDCState, ingress bool) ([]core.Origin, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryDuration)
			defer cancel()
			e, status := host1.RedirectStateChanges(ctx, origin, state, ingress)
			if !status.OK() {
				h.Handle(status)
			}
			return e, status
		},
	}
}()
