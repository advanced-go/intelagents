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

func SetThreshold(h core.ErrorHandler, origin core.Origin, t *threshold1.Entry, observe *Events) {
	if t == nil {
		return
	}
	e, status := observe.IngressThreshold(h, origin)
	if status.OK() {
		t.Percent = e[0].Percent
		t.Value = e[0].Value
		t.Minimum = e[0].Minimum
	} else {
		t.Percent = 95  // percentile
		t.Value = 30000 // milliseconds
		t.Minimum = 0   // no minimum
	}
}

func SetRateLimitingAction(h core.ErrorHandler, origin core.Origin, rl *action1.RateLimiting, exp *Experience) {
	if rl == nil {
		return
	}
	act, status := exp.GetRateLimitingAction(h, origin)
	if status.OK() {
		*rl = act
	} else {
		rl.Limit = defaultLimit
		rl.Burst = defaultBurst
	}
}

func AddExperience(h core.ErrorHandler, origin core.Origin, inf *inference1.Entry, action *action1.RateLimiting, exp *Experience) *core.Status {
	id, status := exp.AddIngressInference(h, origin, *inf)
	if status.OK() {
		action.InferenceId = id
		status = exp.AddRateLimitingAction(h, origin, *action)
	}
	return status
}
