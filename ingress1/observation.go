package ingress1

import (
	"context"
	"github.com/advanced-go/observation/access1"
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
}

func newObservation() *observation {
	return &observation{
		access: func(origin core.Origin) ([]access1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryAccessDuration)
			defer cancel()
			return access1.IngressQuery(ctx, origin)
		},
		inference: func(origin core.Origin) ([]inference1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryInferenceDuration)
			defer cancel()
			return inference1.IngressQuery(ctx, origin)
		},
		addInference: func(e inference1.Entry) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), insertInferenceDuration)
			defer cancel()
			return inference1.IngressInsert(ctx, nil, e)
		},
	}
}
