package egress1

import (
	"context"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	queryTimeseriesDuration = time.Second * 2
	queryAccessDuration     = time.Second * 2
	//queryInferenceDuration  = time.Second * 2

)

// A nod to Linus Torvalds and plain C
type observation struct {
	timeseries   func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status)
	rateLimiting func(h core.ErrorHandler, origin core.Origin) (access1.Entry, *core.Status)
}

func newObservation() *observation {
	return &observation{
		timeseries: func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryTimeseriesDuration)
			defer cancel()
			e, status := timeseries1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},

		rateLimiting: func(h core.ErrorHandler, origin core.Origin) (access1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryAccessDuration)
			defer cancel()
			e, status := access1.IngressRateLimitingQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
	}
}
