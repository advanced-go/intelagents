package ingress1

import (
	"context"
	"fmt"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/stdlib/core"
	fmt2 "github.com/advanced-go/stdlib/fmt"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

// run - ingress controller
func run(c *controller, guid *guidance, observe *observation) {
	if c == nil || guid == nil || observe == nil {
		return
	}
	var prev []access1.Entry
	c.startTicker(0)

	for {
		select {
		case <-c.ticker.C:
			if !guid.shouldProcess(c.origin, c.opsAgent) {
				continue
			}
			testLog(nil, c.uri, "tick")
			curr, status := observe.access(c.origin)
			if !status.OK() && !status.NotFound() {
				c.opsAgent.Handle(status, c.uri)
				continue
			}
			status = processInference(curr, guid.percentile(c.origin, c.opsAgent), observe)
			if !status.OK() {
				c.opsAgent.Handle(status, c.uri)
			}
			prev = curr
			updateTicker(c, prev, curr, observe)
		case msg, open := <-c.ctrlC:
			if !open {
				c.stopTicker()
				return
			}
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(c.ctrlC)
				c.stopTicker()
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
