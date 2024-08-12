package ingress1

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

func resiliencyFunc(r *resiliency, percentile resiliency1.Percentile, observe *observation, exp *experience) ([]timeseries1.Entry, *core.Status) {
	r.handler.AddActivity(r.agentId, "onTick")
	ts, status1 := observe.timeseries(r.handler, r.origin)
	if !status1.OK() || status1.NotFound() {
		return ts, status1
	}
	i := inference1.Entry{}
	/*
		i, status := exp.processInference(c, ts, percentile)
		if !status.OK() {
			return ts, status
		}
	*/
	status := exp.addInference(r.handler, r.origin, i)
	if !status.OK() {
		return ts, status
	}
	action := action1.RateLimiting{InferenceId: 1}
	/*
		actions, status2 := exp.processAction(c, i)
		if !status2.OK() {
			return ts, status2
		}

	*/
	status = exp.addRateLimitingAction(r.handler, r.origin, action)
	return ts, status
}

func resiliencyInitFunc(r *resiliency, observe *observation) *core.Status {

	return core.StatusOK()
}
