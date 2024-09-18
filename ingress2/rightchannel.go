package ingress2

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/intelagents/common2"
	"github.com/advanced-go/stdlib/messaging"
)

const (
	defaultLimit = -1
	defaultBurst = -1
)

// run - ingress resiliency for the RHC
func runResiliencyRHC(r *resiliency, exp *common2.Experience) {
	rateLimiting := action1.RateLimiting{}
	setRateLimiting(r, &rateLimiting, exp)

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
				observe, ok := msg.Body.(*observation)
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
				addExperience(r, inf, action, exp)
			default:
				r.handler.Handle(common.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}
