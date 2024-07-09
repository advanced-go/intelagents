package ingress1

import (
	"context"
	"fmt"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	fmt2 "github.com/advanced-go/stdlib/fmt"
	"github.com/advanced-go/stdlib/messaging"
	"net/http"
	"time"
)

// type logFunc func(ctx context.Context, agentId string, content any) *core.Status
type queryAccessFunc func(ctx context.Context, origin core.Origin) ([]access1.Entry, *core.Status)
type queryInferenceFunc func(ctx context.Context, origin core.Origin) ([]inference1.Entry, *core.Status)
type getGuidanceFunc func()
type insertInferenceFunc func(ctx context.Context, h http.Header, e inference1.Entry) *core.Status

// run - ingress controller
func run(c *controller, access queryAccessFunc, inference queryInferenceFunc, guidance getGuidanceFunc, insert insertInferenceFunc) {
	if c == nil {
		return
	}
	var prev []access1.Entry
	c.StartTicker(0)

	for {
		select {
		case <-c.ticker.C:
			testLog(nil, c.uri, "tick")
			curr, status := access(nil, c.origin)

			if !status.OK() && !status.NotFound() {
				c.handler.Message(messaging.NewStatusMessage(c.handler.Uri(), c.uri, status))
			} else {
				status = processInference(curr, inference, guidance, insert)
				if !status.OK() {
					c.handler.Message(messaging.NewStatusMessage(c.handler.Uri(), c.uri, status))
				}
				prev = curr
			}
			// Update ticker based on changes in RPS
			if updateTicker(c, prev, curr) {
				
			}
		case msg, open := <-c.ctrlC:
			if !open {
				c.StopTicker()
				return
			}
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(c.ctrlC)
				c.StopTicker()
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

func processInference(curr []access1.Entry, inference queryInferenceFunc, guidance getGuidanceFunc, insert insertInferenceFunc) *core.Status {
	e, status := infer(curr, inference, guidance)
	if !status.OK() {
		return status
	}
	status = insert(nil, nil, e)
	if !status.OK() {
		return status
	}
	return core.StatusOK()
}

func updateTicker(c *controller, prev, curr []access1.Entry) bool {
	return false
}
