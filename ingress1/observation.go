package ingress1

import (
	"context"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/experience1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
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
	inference    func(origin core.Origin) ([]inference1.Entry, *core.Status)
	addInference func(e inference1.Entry) *core.Status
	experience   func(origin core.Origin) ([]experience1.Entry, *core.Status)
}

func newObservation(handler func(status *core.Status, _ string) *core.Status) *observation {
	if handler == nil {
		handler = func(status *core.Status, _ string) *core.Status {
			return status
		}
	}
	return &observation{
		access: func(origin core.Origin) ([]access1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryAccessDuration)
			defer cancel()
			e, status := access1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				handler(status, "")
			}
			return e, status
		},
		inference: func(origin core.Origin) ([]inference1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryInferenceDuration)
			defer cancel()
			e, status := inference1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				handler(status, "")
			}
			return e, status
		},
		addInference: func(e inference1.Entry) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), insertInferenceDuration)
			defer cancel()
			status := inference1.IngressInsert(ctx, nil, e)
			if !status.OK() && !status.NotFound() {
				handler(status, "")
			}
			return status
		},
		experience: func(origin core.Origin) ([]experience1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryInferenceDuration)
			defer cancel()
			e, status := experience1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				handler(status, "")
			}
			return e, status
		},
	}
}
