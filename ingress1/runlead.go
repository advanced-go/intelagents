package ingress1

import (
	"github.com/advanced-go/stdlib/messaging"
)

func runLead(a *lead, observe *observation, guide *guidance, ops *operations) {
	if a == nil || observe == nil || guide == nil || ops == nil {
		return
	}

	for {
		select {
		case msg := <-a.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				a.shutdown()
				return
			case messaging.RestartEvent:
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
