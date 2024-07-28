package ingress1

import (
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

func runRedirect(r *redirect, observe *observation, guide *guidance, ops *operations) {
	//observe == nil || guide == nil || ops == nil {
	if r == nil {
		return
	}
	tick := time.Tick(r.interval)

	for {
		select {
		case <-tick:

		// control channel
		case msg := <-r.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.shutdown()
				return
			// How to handle duplicate restart events as multiple pods will be sending restarts?
			//
			case messaging.RestartEvent:
				// TODO : read/update from guidance
				//m := messaging.NewControlMessage(a.dependencyAgent.Uri(),a.uri,messaging.RestartEvent)
				//a.dependencyAgent.Message(m)
				//a.controllers.Broadcast(m)

			// Duplicates?? Should not happen.
			// Changesets only contain the changes
			case messaging.ChangesetApplyEvent:
			case messaging.ChangesetRollbackEvent:
				// TODO : apply and rollback changeset
			default:
			}
		default:
		}
	}
}
