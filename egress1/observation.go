package egress1

import (
	"context"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/experience1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	queryAccessDuration     = time.Second * 2
	queryInferenceDuration  = time.Second * 2
	insertInferenceDuration = time.Second * 2
)

// A nod to Linus Torvalds and plain C
type observation struct {
	access       func(origin core.Origin) ([]access1.Entry, *core.Status)
	routing      func(origin core.Origin) ([]access1.Routing, *core.Status)
	inference    func(origin core.Origin) ([]inference1.Entry, *core.Status)
	addInference func(e inference1.Entry) *core.Status
	experience   func(origin core.Origin) ([]experience1.Entry, *core.Status)
}

func newObservation(agent messaging.OpsAgent) *observation {
	return &observation{
		access: func(origin core.Origin) ([]access1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryAccessDuration)
			defer cancel()
			e, status := access1.EgressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				agent.Handle(status, "")
			}
			return e, status
		},
		routing: func(origin core.Origin) ([]access1.Routing, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryAccessDuration)
			defer cancel()
			r, status := access1.EgressRoutingQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				agent.Handle(status, "")
			}
			return r, status
		},
		inference: func(origin core.Origin) ([]inference1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryInferenceDuration)
			defer cancel()
			e, status := inference1.EgressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				agent.Handle(status, "")
			}
			return e, status
		},
		addInference: func(e inference1.Entry) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), insertInferenceDuration)
			defer cancel()
			status := inference1.EgressInsert(ctx, nil, e)
			if !status.OK() {
				agent.Handle(status, "")
			}
			return status
		},
		experience: func(origin core.Origin) ([]experience1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryInferenceDuration)
			defer cancel()
			e, status := experience1.EgressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				agent.Handle(status, "")
			}
			return e, status
		},
	}
}
