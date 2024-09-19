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
	QueryTimeseries func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status)
	QueryLog        func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status)
	GetThreshold    func(h core.ErrorHandler, origin core.Origin) ([]threshold1.Entry, *core.Status)
	GetProfile      func(h core.ErrorHandler) (*threshold1.Profile, *core.Status)
}

var IngressEvents = func() *Events {
	return &Events{
		QueryTimeseries: func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.QueryIngress(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		QueryLog: func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), logDuration)
			defer cancel()
			e, status := log1.QueryIngress(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		GetThreshold: func(h core.ErrorHandler, origin core.Origin) ([]threshold1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), thresholdDuration)
			defer cancel()
			e, status := threshold1.GetIngress(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		GetProfile: func(h core.ErrorHandler) (*threshold1.Profile, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), thresholdDuration)
			defer cancel()
			e, status := threshold1.GetProfile(ctx)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
	}
}()

var EgressEvents = func() *Events {
	return &Events{
		QueryTimeseries: func(h core.ErrorHandler, origin core.Origin) ([]timeseries1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.QueryEgress(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		QueryLog: func(h core.ErrorHandler, origin core.Origin) ([]log1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := log1.QueryEgress(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		GetThreshold: func(h core.ErrorHandler, origin core.Origin) ([]threshold1.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), thresholdDuration)
			defer cancel()
			e, status := threshold1.GetEgress(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		GetProfile: func(h core.ErrorHandler) (*threshold1.Profile, *core.Status) {
			return IngressEvents.GetProfile(h)
			/*ctx, cancel := context.WithTimeout(context.Background(), thresholdDuration)
			defer cancel()
			e, status := threshold1.GetProfile(ctx)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status

			*/
		},
	}
}()
