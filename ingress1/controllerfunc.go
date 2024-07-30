package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/stdlib/core"
)

func controllerFunc(c *controller, percentile percentile1.Entry, observe *observation, exp *experience) ([]access1.Entry, *core.Status) {
	c.handler.AddActivity(c.agentId, "onTick")
	curr, status1 := observe.access(c.handler, c.origin)
	if !status1.OK() || status1.NotFound() {
		return curr, status1
	}
	i, status := exp.processInference(c, curr, percentile)
	if !status.OK() {
		return curr, status
	}
	status = exp.addInference(c.handler, c.origin, i)
	if !status.OK() {
		return curr, status
	}
	actions, status2 := exp.processAction(c, i)
	if !status2.OK() {
		return curr, status2
	}
	status = exp.addAction(c.handler, c.origin, actions)
	return curr, status
}

/*
func processControlInference(c *controller, e []access1.Entry, percentile percentile1.Entry, observe *observation, exp *experience, inf *inference, ops *operations) (inference1.Entry, *core.Status) {
	i, status := inf.process(c, e, percentile, exp, ops)
	if !status.OK() {
		return inference1.Entry{}, status
	}
	status = exp.addInference(i)
	return i, status
}

func processControlAction(c *controller, i inference1.Entry, exp *experience, act *action, ops *operations) *core.Status {
	actions, status := act.process(i, ops)
	if !status.OK() {
		return status
	}
	status = act.insert(c.origin, actions)
	return status
}


*/
