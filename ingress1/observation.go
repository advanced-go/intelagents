package ingress1

import (
	"context"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	queryAccessDuration = time.Second * 2
)

// A nod to Linus Torvalds and plain C
type observation struct {
	access         func(h core.ErrorHandler, origin core.Origin) ([]access1.Entry, *core.Status)
	accessRedirect func(h core.ErrorHandler, origin core.Origin) ([]access1.Entry, *core.Status)
}

var observe = func() *observation {
	return &observation{
		access: func(h core.ErrorHandler, origin core.Origin) ([]access1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryAccessDuration)
			defer cancel()
			e, status := access1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		accessRedirect: func(h core.ErrorHandler, origin core.Origin) ([]access1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryAccessDuration)
			defer cancel()
			e, status := access1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
	}
}()
