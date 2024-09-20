package common1

import (
	"context"
	"github.com/advanced-go/events/timeseries1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	timeseriesDuration = time.Second * 2
)

// Events - events interface, with a nod to Linus Torvalds and plain C
type Events struct {
	PercentThresholdSLO      func(h core.ErrorHandler, origin core.Origin) (timeseries1.Threshold, *core.Status)
	PercentThresholdQuery    func(h core.ErrorHandler, origin core.Origin, from time.Time, to time.Time) (timeseries1.Threshold, *core.Status)
	StatusCodeThresholdQuery func(h core.ErrorHandler, origin core.Origin, from time.Time, to time.Time, statusCodes string) (timeseries1.Threshold, *core.Status)
	GetProfile               func(h core.ErrorHandler) (*timeseries1.Profile, *core.Status)
}

var TimeseriesEvents = func() *Events {
	return &Events{
		PercentThresholdSLO: func(h core.ErrorHandler, origin core.Origin) (timeseries1.Threshold, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.PercentileThresholdSLO(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		PercentThresholdQuery: func(h core.ErrorHandler, origin core.Origin, from time.Time, to time.Time) (timeseries1.Threshold, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.PercentileThresholdQuery(ctx, origin, from, to)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		StatusCodeThresholdQuery: func(h core.ErrorHandler, origin core.Origin, from time.Time, to time.Time, statusCodes string) (timeseries1.Threshold, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.StatusCodeThresholdQuery(ctx, origin, from, to, statusCodes)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
		GetProfile: func(h core.ErrorHandler) (*timeseries1.Profile, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries1.GetProfile(ctx)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
	}
}()
