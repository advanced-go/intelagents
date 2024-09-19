package redirect1

import (
	"github.com/advanced-go/events/threshold1"
	"github.com/advanced-go/intelagents/common2"
	"github.com/advanced-go/stdlib/messaging"
)

// run - ingress resiliency for the LHC
// r.handler.AddActivity(r.agentId, "onTick")
func runRedirectLHC(r *redirect, observe *common2.Events) {
	origin := redirectOrigin(r.origin, r.state.Location)
	ticker := messaging.NewTicker(redirectDuration)
	// Set the threshold to the current host's, and use that to compare to the redirect host's actual threshold
	limit := threshold1.Entry{}
	common2.SetPercentileThreshold(r.handler, r.origin, &limit, observe)

	ticker.Start(-1)
	for {
		// observation processing
		select {
		case <-ticker.C():
			actual, status := observe.GetThreshold(r.handler, origin)
			if status.OK() {
				m := messaging.NewRightChannelMessage("", r.agentId, messaging.ObservationEvent, common2.NewObservation(actual[0], limit))
				r.Message(m)
			}
		default:
		}
		// message processing
		select {
		case msg := <-r.lhc.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				ticker.Stop()
				r.lhc.Close()
				return
			//case messaging.DataChangeEvent:
			default:
				r.handler.Handle(common2.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}
