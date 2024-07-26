package ingress1

import (
	"github.com/advanced-go/stdlib/messaging"
)

func runLead(a *lead, observe *observation, guide *guidance, ops *operations) {
	if a == nil || observe == nil || guide == nil || ops == nil {
		return
	}
	entry, status := guide.controllers(a.origin)
	if entry.EntryId != 0 || status != nil {
	}
	for {
		select {
		case msg := <-a.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				a.shutdown()
				return
			case messaging.HostStartupEvent:
				//if
				a.controller.Message(msg)
			case messaging.ChangesetApplyEvent:
				a.controller.Message(msg)
			case messaging.ChangesetRollbackEvent:
				a.controller.Message(msg)
			default:
			}
		default:
		}
	}
}
