package common

import (
	"context"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	timeseriesDuration = time.Second * 2
)

// A nod to Linus Torvalds and plain C

// Observation - observation functions
type Observation struct {
	IngressTimeseries func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status)
	EgressTimeseries  func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status)
}

var Observe = func() *Observation {
	return &Observation{
		IngressTimeseries: func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
		EgressTimeseries: func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.EgressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return e, status
		},
	}
}()
