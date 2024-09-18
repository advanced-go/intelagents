package ingress2

import (
	"github.com/advanced-go/events/threshold1"
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/intelagents/common2"
	"github.com/advanced-go/stdlib/core"
)

func setThreshold(r *resiliency, t *threshold1.Entry, observe *common2.Events) {
	if r == nil || t == nil {
		return
	}
	e, status := observe.Threshold(r.handler, r.origin)
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

func setRateLimiting(r *resiliency, rl *action1.RateLimiting, exp *common2.Experience) {
	if r == nil || rl == nil {
		return
	}
	act, status := exp.GetRateLimitingAction(r.handler, r.origin)
	if status.OK() {
		*rl = act
	} else {
		rl.Limit = defaultLimit
		rl.Burst = defaultBurst
	}
}

func addExperience(r *resiliency, inf *inference1.Entry, action *action1.RateLimiting, exp *common2.Experience) *core.Status {
	id, status := exp.AddIngressInference(r.handler, r.origin, *inf)
	if status.OK() {
		action.InferenceId = id
		status = exp.AddRateLimitingAction(r.handler, r.origin, *action)
	}
	return status
}
