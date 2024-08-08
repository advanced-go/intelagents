package ingress1

import (
	"github.com/advanced-go/stdlib/core"
)

func redirectFunc(r *redirect, observe *observation, exp *experience) *core.Status {
	r.handler.AddActivity(r.agentId, "onTick")
	/*
		ts, status1 := observe.timeseries(c.handler, c.origin)
		if !status1.OK() || status1.NotFound() {
			return ts, status1
		}
		i, status := exp.processInference(c, ts, percentile)
		if !status.OK() {
			return ts, status
		}
		status = exp.addInference(c.handler, c.origin, i)
		if !status.OK() {
			return ts, status
		}
		actions, status2 := exp.processAction(c, i)
		if !status2.OK() {
			return ts, status2
		}
		status = exp.addAction(c.handler, c.origin, actions)

	*/
	return core.StatusOK()
}

func redirectInitFunc(r *redirect, observe *observation) *core.Status {
	entry, status := observe.redirect(r.handler, r.origin)
	if status.OK() {
		r.state.host = r.origin.Host
		r.state.location = entry.From
		r.state.percent = 0
		return status
	}
	if !status.NotFound() {
		r.handler.Handle(status, "")
	}
	return status
}
