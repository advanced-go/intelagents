package common1

import (
	"github.com/advanced-go/events/timeseries1"
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	DefaultStatusCodes = "5xx"
)

type Observation struct {
	Actual timeseries1.Threshold
	Limit  timeseries1.Threshold
}

func NewObservation(actual, limit timeseries1.Threshold) *Observation {
	o := new(Observation)
	o.Actual = actual
	o.Limit = limit
	return o
}

func SetPercentileThreshold(h core.ErrorHandler, origin core.Origin, t *timeseries1.Threshold, observe *Events) {
	if t == nil {
		return
	}
	e, status := observe.PercentThresholdSLO(h, origin)
	if status.OK() {
		t.Percent = e.Percent
		t.Value = e.Value
		t.Minimum = e.Minimum
	} else {
		timeseries1.InitPercentileThreshold(t)
	}
}

func SetStatusCodesThreshold(h core.ErrorHandler, origin core.Origin, t *timeseries1.Threshold, from, to time.Time, statusCodes string, observe *Events) {
	if t == nil {
		return
	}
	if statusCodes == "" {
		statusCodes = DefaultStatusCodes
	}
	e, status := observe.StatusCodeThresholdQuery(h, origin, from, to, statusCodes)
	if status.OK() {
		t.Percent = e.Percent
		t.Value = e.Value
		t.Minimum = e.Minimum
	} else {
		timeseries1.InitStatusCodeThreshold(t)
	}
}

func SetRateLimitingAction(h core.ErrorHandler, origin core.Origin, a *action1.RateLimiting, exp *Experience) {
	if a == nil {
		return
	}
	act, status := exp.ActiveRateLimitingAction(h, origin)
	if status.OK() {
		*a = act
	} else {
		action1.InitRateLimiting(a)
	}
}

func SetRoutingAction(h core.ErrorHandler, origin core.Origin, a *action1.Routing, exp *Experience) {
	if a == nil {
		return
	}
	act, status := exp.ActiveRoutingAction(h, origin)
	if status.OK() {
		*a = act
	} else {
		action1.InitRouting(a)
	}
}

func AddRateLimitingExperience(h core.ErrorHandler, origin core.Origin, inf *inference1.Entry, a *action1.RateLimiting, exp *Experience) *core.Status {
	id, status := exp.AddInference(h, origin, inf)
	if status.OK() {
		a.InferenceId = id
		status = exp.AddRateLimitingAction(h, origin, a)
	}
	return status
}

func AddRoutingExperience(h core.ErrorHandler, origin core.Origin, inf *inference1.Entry, a *action1.Routing, exp *Experience) *core.Status {
	id, status := exp.AddInference(h, origin, inf)
	if status.OK() {
		a.InferenceId = id
		status = exp.AddRoutingAction(h, origin, a)
	}
	return status
}
