package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
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
func runControl(c *controller, observe *observation, exp *experience, guide *guidance, infer *inference, ops *operations) {
	if c == nil || observe == nil || exp == nil || guide == nil || infer == nil || ops == nil {
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
			status = infer.process(curr, percentile, observe)
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
			infer.updateTicker(c, exp)
		}
	}
}
