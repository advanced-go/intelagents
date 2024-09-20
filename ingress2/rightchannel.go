package ingress2

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/intelagents/common1"
	"github.com/advanced-go/stdlib/messaging"
)

// run - ingress resiliency for the RHC
func runResiliencyRHC(r *resiliency, exp *common1.Experience) {
	rateLimiting := action1.NewRateLimiting()
	common1.SetRateLimitingAction(r.handler, r.origin, rateLimiting, exp)

	for {
		// message processing
		select {
		case msg := <-r.rhc.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.rhc.Close()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			case messaging.ObservationEvent:
				r.handler.AddActivity(r.agentId, messaging.ObservationEvent)
				observe, ok := msg.Body.(*common1.Observation)
				if !ok {
					continue
				}
				inf := runInference(r, observe)
				if inf == nil {
					continue
				}
				action := newAction(inf)
				rateLimiting.Limit = action.Limit
				rateLimiting.Burst = action.Burst
				common1.AddRateLimitingExperience(r.handler, r.origin, inf, action, exp)
			default:
				r.handler.Handle(common1.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}
