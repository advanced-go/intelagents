package common2

import (
	"context"
	"github.com/advanced-go/access/log1"
	"github.com/advanced-go/access/threshold1"
	"github.com/advanced-go/access/timeseries1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	timeseriesDuration = time.Second * 2
)

// Access - access functions struct, a nod to Linus Torvalds and plain C
type Access struct {
	IngressTimeseries func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status)
	EgressTimeseries  func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status)

	IngressLog func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status)
	EgressLog  func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status)

	Threshold func(h core.ErrorHandler, origin core.Origin) ([]threshold1.Entry, *core.Status)
}

var Observ = func() *Access {
	return &Access{
		IngressTimeseries: func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		EgressTimeseries: func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.EgressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		IngressLog: func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := log1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		EgressLog: func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := log1.EgressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		Threshold: func(h core.ErrorHandler, origin core.Origin) ([]threshold1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := threshold1.Query(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
	}
}()
