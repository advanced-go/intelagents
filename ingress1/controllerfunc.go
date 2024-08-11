package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

func controllerFunc(c *controller, percentile percentile1.Entry, observe *observation, exp *experience) ([]timeseries1.Entry, *core.Status) {
	c.handler.AddActivity(c.agentId, "onTick")
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
	/*
		actions, status2 := exp.processAction(c, i)
		if !status2.OK() {
			return ts, status2
		}
		status = exp.addAction(c.handler, c.origin, actions)

	*/
	return ts, status
}

func controllerInitFunc(c *controller, observe *observation) *core.Status {
	entry, status := observe.rateLimiting(c.handler, c.origin)
	if status.OK() {
		c.state.rateLimit = entry.RateLimit
		c.state.rateBurst = int(entry.RateBurst)
		return status
	}
	if !status.NotFound() {
		c.handler.Handle(status, "")
	}
	return status
}
