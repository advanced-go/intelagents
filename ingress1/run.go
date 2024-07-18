package ingress1

import (
	"context"
	"fmt"
	"github.com/advanced-go/intelagents/guidance1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/percentile1"
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
func run(c *controller, observe *observation) {
	if c == nil || observe == nil {
		return
	}
	var prev []access1.Entry
	percentile, status := observe.percentile(percentileDuration, defaultPercentile, c.origin)
	if !status.OK() {
		c.handler.Handle(status, "")
	}
	c.ticker.Start(0)
	c.poller.Start(0)
	for {
		select {
		// main : on tick -> observe access -> process inference with percentile -> create action
		case <-c.ticker.C():
			if !guidance1.IsScheduled() {
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
			percentile, status = observe.percentile(percentileDuration, percentile, c.origin)
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
