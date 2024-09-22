package redirect1

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/guidance/redirect1"
	"github.com/advanced-go/intelagents/common1"
	"github.com/advanced-go/stdlib/messaging"
)

func runRedirectRHC(r *redirect, exp *common1.Experience, guide *common1.RedirectGuidance) {
	// Set existing routing action state
	routing := action1.NewRouting()
	common1.SetRoutingAction(r.handler, r.origin, routing, exp)

	// Set initial step percentage from current/latest routing action
	stepPercent := 0
	if routing.Percentage > 0 {
		stepPercent = routing.Percentage
	}
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
				if obs := castObservation(r.handler, r.agentId, msg); obs != nil {
					// Compare actual and limit for percentile and status code thresholds
					status := runInference(r, obs)
					if !status.OK() {
						// Add failure status to redirect, adding reason for failure in comment, and exit agent
						common1.RedirectGuide.AddIngressStatus(r.handler, r.origin, redirect1.RedirectStatusFailed, "")
						// Add nil routing action to close the redirect
						status = exp.AddRoutingAction(r.handler, r.origin, nil)
						r.rhc.Close()
						return
					}
					// Add status for passing step
					common1.RedirectGuide.AddIngressStatus(r.handler, r.origin, redirect1.RedirectStatusUpdate, "updated step percentage")

					// Update step
					if stepPercent == 100 {
						common1.RedirectGuide.AddIngressStatus(r.handler, r.origin, redirect1.RedirectStatusSucceeded, "updated step percentage")
						// Add a nil routing action so the routing state is complete
						status = exp.AddRoutingAction(r.handler, r.origin, nil)
						r.rhc.Close()
						return
					}
					stepPercent = updatePercentage(stepPercent)

					// Add new routing action
					action := action1.NewRouting()
					action.AgentId = r.agentId
					status = exp.AddRoutingAction(r.handler, r.origin, action)
				}
			default:
			}
		default:
		}
	}
}
