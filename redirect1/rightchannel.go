package redirect1

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/intelagents/common2"
	"github.com/advanced-go/stdlib/messaging"
)

func runRedirectRHC(r *redirect, exp *common2.Experience, guide *common2.Guidance) {
	routing := action1.NewRouting()
	common2.SetRoutingAction(r.handler, r.origin, routing, exp)

	// Startup: if the current routing is active, then initialize the state percentage
	//          if the current routing is not active, then send a new action to start redirecting
	for {
		select {
		case msg := <-r.rhc.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.rhc.Close()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			case messaging.ObservationEvent:
				r.handler.AddActivity(r.agentId, messaging.ObservationEvent)
				_, ok := msg.Body.(*common2.Observation)
				if !ok {
					continue
				}
				// if the actual meets limit, then look at retries
				/*
					inf := runInference(r, observe)
					if inf == nil {
						continue
					}
					action := newAction(inf)
					rateLimiting.Limit = action.Limit
					rateLimiting.Burst = action.Burst
					common2.AddRateLimitingExperience(r.handler, r.origin, inf, action, exp)

				*/
			default:
			}
		default:
		}
	}
}
