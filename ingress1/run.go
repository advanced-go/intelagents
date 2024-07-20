package ingress1

import (
	"context"
	"fmt"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/guidance/schedule1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/stdlib/core"
	fmt2 "github.com/advanced-go/stdlib/fmt"
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
func run(c *controller, observe *observation, guide *guidance) {
	if c == nil || observe == nil || guide == nil {
		return
	}
	var prev []access1.Entry

	percentile, status := guide.percentile(percentileDuration, defaultPercentile, c.origin, c.handler)
	c.ticker.Start(0)
	c.poller.Start(0)
	for {
		select {
		// main : on tick -> observe access -> process inference with percentile -> create action
		case <-c.ticker.C():
			if !schedule1.IsIngressControllerScheduled(c.origin) {
				continue
			}
			testLog(nil, c.uri, "tick")
			curr, status1 := observe.access(c.origin)
			if !status1.OK() {
				if !status1.NotFound() {
					c.handler.Handle(status1, c.uri)
				}
				continue
			}
			status1 = processInference(curr, percentile, observe)
			if !status1.OK() {
				c.handler.Handle(status, c.uri)
				continue
			}
			prev = curr
			updateTicker(c, prev, curr, observe)
		// poll : update percentile
		case <-c.poller.C():
			percentile, status = guide.percentile(percentileDuration, percentile, c.origin, c.handler)
			if !status.OK() {
				c.handler.Handle(status, "")
			}
		// control channel
		case msg := <-c.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(c.ctrlC)
				c.stopTickers()
				testLog(nil, c.uri, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}

func testLog(_ context.Context, agentId string, content any) *core.Status {
	fmt.Printf("test: activity1.Log() -> %v : %v : %v\n", fmt2.FmtRFC3339Millis(time.Now().UTC()), agentId, content)
	return core.StatusOK()
}

func processInference(curr []access1.Entry, percentile percentile1.Entry, observe *observation) *core.Status {
	e, status := infer(curr, percentile, observe)
	if !status.OK() {
		return status
	}
	return observe.addInference(e)
}

func updateTicker(c *controller, prev, curr []access1.Entry, observe *observation) *core.Status {
	// Need to insert another inference entry for changing the ticker based on RPS
	return core.StatusOK()
}
