package dependency1

import (
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

func run(a *dependency) {
	if a == nil {
		return
	}
	tick := time.Tick(a.interval)

	// TODO: Read and update from configurations
	for {
		select {
		case <-tick:

		// control channel
		case msg := <-a.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				return
			case messaging.RestartEvent:
				// TODO : restart
			case messaging.ChangesetApplyEvent:
			case messaging.ChangesetRollbackEvent:
				// TODO : apply and rollback
			default:
			}
		default:
		}
	}
}
