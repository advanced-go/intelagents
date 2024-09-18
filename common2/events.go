package common2

import (
	"context"
	"github.com/advanced-go/events/log1"
	"github.com/advanced-go/events/threshold1"
	"github.com/advanced-go/events/timeseries1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	timeseriesDuration = time.Second * 2
	logDuration        = time.Second * 2
	thresholdDuration  = time.Second * 2
)

// Events - access functions struct, a nod to Linus Torvalds and plain C
type Events struct {
	IngressTimeseries func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status)
	EgressTimeseries  func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status)

	IngressLog func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status)
	EgressLog  func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status)

	IngressThreshold func(h core.ErrorHandler, origin core.Origin) ([]threshold1.Entry, *core.Status)
	EgressThreshold  func(h core.ErrorHandler, origin core.Origin) ([]threshold1.Entry, *core.Status)

	GetProfile func(h core.ErrorHandler) (*threshold1.Profile, *core.Status)
}

var Event = func() *Events {
	return &Events{
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
			ctx, cancel := context.WithTimeout(context.Background(), logDuration)
			defer cancel()
			e, status := log1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		EgressLog: func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), logDuration)
			defer cancel()
			e, status := log1.EgressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		IngressThreshold: func(h core.ErrorHandler, origin core.Origin) ([]threshold1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), thresholdDuration)
			defer cancel()
			e, status := threshold1.IngressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		EgressThreshold: func(h core.ErrorHandler, origin core.Origin) ([]threshold1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), thresholdDuration)
			defer cancel()
			e, status := threshold1.EgressQuery(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		GetProfile: func(h core.ErrorHandler) (*threshold1.Profile, *core.Status) {
			//ctx, cancel := context.WithTimeout(context.Background(), thresholdDuration)
			//defer cancel()
			e, status := threshold1.GetProfile()
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
	}
}()
