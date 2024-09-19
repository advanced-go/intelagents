package common2

import (
	"github.com/advanced-go/events/threshold1"
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/stdlib/core"
)

const (
	defaultLimit = -1
	defaultBurst = -1
)

type Observation struct {
	Actual threshold1.Entry
	Limit  threshold1.Entry
}

func NewObservation(actual, limit threshold1.Entry) *Observation {
	o := new(Observation)
	o.Actual = actual
	o.Limit = limit
	return o
}

func SetPercentileThreshold(h core.ErrorHandler, origin core.Origin, t *threshold1.Entry, observe *Events) {
	if t == nil {
		return
	}
	e, status := observe.GetThreshold(h, origin)
	if status.OK() {
		t.Percent = e[0].Percent
		t.Value = e[0].Value
		t.Minimum = e[0].Minimum
	} else {
		threshold1.InitPercentileThreshold(t)
	}
}

func SetStatusCodesThreshold(h core.ErrorHandler, origin core.Origin, t *threshold1.Entry, observe *Events) {
	if t == nil {
		return
	}
	e, status := observe.GetThreshold(h, origin)
	if status.OK() {
		t.Percent = e[0].Percent
		t.Value = e[0].Value
		t.Minimum = e[0].Minimum
	} else {
		threshold1.InitStatusCodeThreshold(t)
	}
}

func SetRateLimitingAction(h core.ErrorHandler, origin core.Origin, a *action1.RateLimiting, exp *Experience) {
	if a == nil {
		return
	}
	act, status := exp.GetRateLimitingAction(h, origin)
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
	act, status := exp.GetRoutingAction(h, origin)
	if status.OK() {
		*a = act
	} else {
		action1.InitRouting(a)
	}
}

func AddRateLimitingExperience(h core.ErrorHandler, origin core.Origin, inf *inference1.Entry, a *action1.RateLimiting, exp *Experience) *core.Status {
	id, status := exp.AddInference(h, origin, *inf)
	if status.OK() {
		a.InferenceId = id
		status = exp.AddRateLimitingAction(h, origin, *a)
	}
	return status
}

func AddRoutingExperience(h core.ErrorHandler, origin core.Origin, inf *inference1.Entry, a *action1.Routing, exp *Experience) *core.Status {
	id, status := exp.AddInference(h, origin, *inf)
	if status.OK() {
		a.InferenceId = id
		status = exp.AddRoutingAction(h, origin, *a)
	}
	return status
}
