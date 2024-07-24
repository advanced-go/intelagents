package ingress1

import (
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

func runLead(a *lead, observe *observation, guide *guidance, ops *operations) {
	if a == nil || observe == nil || guide == nil || ops == nil {
		return
	}
	tick := time.Tick(a.interval)

	status := startup(a, messaging.NewControlMessage("", "", messaging.StartupEvent), observe, guide)
	if !status.OK() {
		a.Handle(status, "")
	}
	for {
		select {
		case <-tick:

		// control channel
		case msg := <-a.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				return
			// How to handle duplicate restart events as multiple pods will be sending restarts?
			//
			case messaging.RestartEvent:
				// TODO : read/update from guidance
				//m := messaging.NewControlMessage(a.dependencyAgent.Uri(),a.uri,messaging.RestartEvent)
				//a.dependencyAgent.Message(m)
				//a.controllers.Broadcast(m)
				status = startup(a, msg, observe, guide)
				if !status.OK() {

				}
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

// Need to read the routing information from the access log
func startup(a *lead, msg *messaging.Message, observe *observation, guide *guidance) *core.Status {
	//entries, status := guide.controllers(a.origin)
	//if !status.OK() {
	//	return status
	//}
	//	routing, status1 := observe.routing(a.origin)
	//	if !status1.OK() {
	//		return status1
	//	}
	//	status = startupDependency(a, msg, entries[0])
	//	if !status.OK() {
	//		return status
	//	}
	return core.StatusOK() //startupControllers(a, msg, entries, routing)
}
