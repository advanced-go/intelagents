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
	updateRateLimiting(r, &rateLimiting, exp)

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

				inf, action, status := runInference(r, nil, exp)
			//	r.handler.AddActivity(r.agentId, fmt.Sprintf("%v - %v", msg.Event(), msg.ContentType()))
			//	processDataChangeEvent(r, msg, guide)
			default:
				r.handler.Handle(common.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}

func updateRateLimiting(r *resiliency, rl *action1.RateLimiting, exp *common2.Experience) {
	if r == nil || rl == nil {
		return
	}
	act, status := exp.GetRateLimitingAction(r.handler, r.origin)
	if status.OK() {
		*rl = act
	} else {
		rl.Limit = defaultLimit
		rl.Burst = defaultBurst
	}
}
