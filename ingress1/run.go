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
	c.StartTicker(0)

	for {
		select {
		case <-c.ticker.C:
			testLog(nil, c.uri, "tick")

			status := processAssignment()
			if !status.OK() && !status.NotFound() {
				c.handler.Message(messaging.NewStatusMessage(c.handler.Uri(), c.uri, status))
			}
			// Update ticker based on changes in RPS
			updateTicker(c)
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

func processAssignment() *core.Status {
	return core.StatusOK()
}

func updateTicker(c *controller) {
	
}
