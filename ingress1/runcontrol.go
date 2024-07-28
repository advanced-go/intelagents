package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	percentileDuration = time.Second * 2
)

var (
	defaultPercentile = percentile1.Entry{Percent: 99, Latency: 2000}
)

// run - ingress controller
func runControl(c *controller, observe *observation, exp *experience, guide *guidance, infer *inference, act *action, ops *operations) {
	if c == nil || observe == nil || exp == nil || guide == nil || infer == nil || act == nil || ops == nil {
		return
	}
	percentile, _ := guide.percentile(percentileDuration, defaultPercentile, c.origin)

	c.ticker.Start(0)
	c.poller.Start(0)

	for {
		select {
		case <-c.ticker.C():
			// main : on tick -> observe access -> process inference with percentile -> create action
			if !guide.isScheduled(c.origin) {
				continue
			}
			ops.addActivity(c.uri, "tick")
			curr, status := observe.access(c.origin)
			if !status.OK() {
				continue
			}
			i, status1 := processControlInference(c, curr, percentile, observe, exp, infer, ops)
			if status1.OK() {
				processControlAction(c, i, exp, act, ops)
			}
		case <-c.poller.C():
			percentile, _ = guide.percentile(percentileDuration, percentile, c.origin)
		case msg := <-c.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				c.shutdown()
				ops.addActivity(c.uri, messaging.ShutdownEvent)
				return
			default:
			}
		default:
			// TODO: should this be scheduled, and what data is needed?
			exp.updateTicker(c, ops)
		}
	}
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
