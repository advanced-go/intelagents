package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
)

func controllerFunc(c *controller, percentile percentile1.Entry, observe *observation, exp *experience, action *action) ([]access1.Entry, *core.Status) {
	c.handler.AddActivity(c.agentId, "onTick")
	curr, status1 := observe.access(c.origin)
	if !status1.OK() && !status1.NotFound() {
		c.handler.Handle(status1, "")
		return curr, status1
	}
	i, status := infer(c, curr, percentile, exp)
	if !status.OK() {
		c.handler.Handle(status1, "")
		return curr, status
	}
	status = observe.addInference(i)
	if !status.OK() {
		c.handler.Handle(status, "")
		return curr, status
	}
	actions, status2 := act(i, c.handler)
	if !status2.OK() {
		c.handler.Handle(status2, "")
		return curr, status2
	}
	status = action.insert(c.origin, actions)
	if !status.OK() {
		c.handler.Handle(status, "")
		//return status
	}
	return curr, status
}

func processControlInference(c *controller, e []access1.Entry, percentile percentile1.Entry, observe *observation, exp *experience, inf *inference, ops *operations) (inference1.Entry, *core.Status) {
	i, status := inf.process(c, e, percentile, exp, ops)
	if !status.OK() {
		return inference1.Entry{}, status
	}
	status = observe.addInference(i)
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
