package guidance1

import (
	"context"
	"fmt"
	"github.com/advanced-go/stdlib/core"
	fmt2 "github.com/advanced-go/stdlib/fmt"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

// run - schedule
func run(c *schedule) {
	if c == nil {
		return
	}
	c.ticker.Start(0)

	for {
		select {
		case <-c.ticker.C():
			// TODO : poll database for global calendar change

		case msg := <-c.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				close(c.ctrlC)
				c.ticker.Stop()
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
