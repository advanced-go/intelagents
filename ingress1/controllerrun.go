package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

var (
	defaultPercentile = percentile1.Entry{Percent: 99, Latency: 2000}
)

// run - ingress controller
func controllerRun1(c *controller, observe *observation, exp *experience, guide *guidance, infer *inference, act *action, ops *operations) {
	//|| observe == nil || exp == nil || guide == nil || infer == nil || act == nil || ops == nil {
	if c == nil {
		return
	}
	percentile, _ := guide.percentile(c.origin, defaultPercentile)
	c.startup()

	for {
		// main agent processing
		select {
		case <-c.ticker.C():
			// main : on tick -> observe access -> process inference with percentile -> create action
			if !guide.isScheduled(c.origin) {
				continue
			}
			ops.addActivity(c.agentId, "onTick")
			curr, status := observe.access(c.origin)
			if !status.OK() {
				continue
			}
			i, status1 := processControlInference(c, curr, percentile, observe, exp, infer, ops)
			if status1.OK() {
				processControlAction(c, i, exp, act, ops)
			}
		default:
		}
		// secondary processing
		select {
		case <-c.poller.C():
			ops.addActivity(c.agentId, "onPoll")
			percentile, _ = guide.percentile(c.origin, percentile)
		case <-c.revise.C():
			ops.addActivity(c.agentId, "onRevise")
			exp.reviseTicker(c.updateTicker, ops)
		case msg := <-c.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				c.shutdown()
				ops.addActivity(c.agentId, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}

func processControlInference1(c *controller, e []access1.Entry, percentile percentile1.Entry, observe *observation, exp *experience, inf *inference, ops *operations) (inference1.Entry, *core.Status) {
	i, status := inf.process(c, e, percentile, exp, ops)
	if !status.OK() {
		return inference1.Entry{}, status
	}
	status = observe.addInference(i)
	return i, status
}

func processControlAction2(c *controller, i inference1.Entry, exp *experience, act *action, ops *operations) *core.Status {
	actions, status := act.process(i, ops)
	if !status.OK() {
		return status
	}
	status = act.insert(c.origin, actions)
	return status
}
