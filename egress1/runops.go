package egress1

import (
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

func runOps(a *operations) {
	if a == nil {
		return
	}
	tick := time.Tick(a.interval)

	// TODO: read/update from guidance
	for {
		select {
		case <-tick:

		// control channel
		case msg := <-a.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				return
			case messaging.RestartEvent:
				// TODO : read/update from guidance
			case messaging.ChangesetApplyEvent:
			case messaging.ChangesetRollbackEvent:
				// TODO : apply and rollback changeset
			default:
			}
		default:
		}
	}
}
