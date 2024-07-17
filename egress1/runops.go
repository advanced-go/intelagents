package egress1

import (
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

func runOps(a *operations) {
	if a == nil {
		return
	}
	tick := time.Tick(a.interval)

	for {
		select {
		case <-tick:
			status := core.StatusOK()
			if !status.OK() && !status.NotFound() {
				a.handler.Message(messaging.NewStatusMessage(a.handler.Uri(), a.uri, status))
			}
		// control channel
		case msg := <-a.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				return
			default:
			}
		default:
		}
	}
}
