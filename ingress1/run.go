package ingress1

import (
	"context"
	"fmt"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/intelagents/guidance"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	fmt2 "github.com/advanced-go/stdlib/fmt"
	"github.com/advanced-go/stdlib/messaging"
	"net/http"
	"time"
)

type enabledFunc func(origin core.Origin) bool

type queryAccessFunc func(origin core.Origin) ([]access1.Entry, *core.Status)
type queryInferenceFunc func(origin core.Origin) ([]inference1.Entry, *core.Status)

type insertInferenceFunc func(ctx context.Context, h http.Header, e inference1.Entry) *core.Status
type insertIntervalFunc func(ctx context.Context, h http.Header, e inference1.Entry) *core.Status

// run - ingress controller
func run(c *controller, enabled enabledFunc, access queryAccessFunc, inference queryInferenceFunc, insert insertInferenceFunc) {
	if c == nil {
		return
	}
	var prev []access1.Entry
	c.startTicker(0)

	for {
		select {
		case <-c.ticker.C:
			if !enabled(c.origin) {
				continue
			}
			percentile := guidance.GetPercentile(c.origin, c.opsAgent)
			testLog(nil, c.uri, "tick")
			curr, status := access(c.origin)
			if !status.OK() && !status.NotFound() {
				c.opsAgent.Handle(status, c.uri)
				continue
			}
			status = processInference(curr, inference, percentile, insert)
			if !status.OK() {
				c.opsAgent.Handle(status, c.uri)
			}
			prev = curr
			updateTicker(c, prev, curr, insert)
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

func processInference(curr []access1.Entry, inference queryInferenceFunc, percentile percentile1.Entry, insert insertInferenceFunc) *core.Status {
	e, status := infer(curr, inference, percentile)
	if !status.OK() {
		return status
	}
	return insert(nil, nil, e)
}

func updateTicker(c *controller, prev, curr []access1.Entry, insert insertInferenceFunc) *core.Status {
	// Need to insert another inference entry for changing the ticker based on RPS
	return core.StatusOK()
}
