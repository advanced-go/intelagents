package ingress1

import (
	"github.com/advanced-go/stdlib/messaging"
)

func leadRun(l *lead, observe *observation, guide *guidance, ops *operations) {
	//} observe == nil || guide == nil || ops == nil {
	if l == nil {
		return
	}
	entry, status := guide.controllers(l.origin)
	if entry.EntryId != 0 || status != nil {
	}
	for {
		select {
		case msg := <-l.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				l.shutdown()
				return
			case messaging.HostStartupEvent:
				//if
				l.controller.Message(msg)
			case messaging.ChangesetApplyEvent:
				l.controller.Message(msg)
			case messaging.ChangesetRollbackEvent:
				l.controller.Message(msg)
			default:
			}
		default:
		}
	}
}
