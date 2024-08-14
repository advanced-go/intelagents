package ingress1

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

// A nod to Linus Torvalds and plain C
type resiliencyFunc struct {
	startup   func(r *resiliency, exp *experience) (*resiliency1.Percentile, *core.Status)
	process   func(r *resiliency, percentile *resiliency1.Percentile, observe *observation, exp *experience) ([]timeseries1.Entry, *core.Status)
	inference func(r *resiliency, entry []timeseries1.Entry, percentile *resiliency1.Percentile) (inference1.Entry, *core.Status)
	action    func(r *resiliency, entry inference1.Entry) (action1.RateLimiting, *core.Status)
}

var resilience = func() *resiliencyFunc {
	return &resiliencyFunc{
		startup: func(r *resiliency, exp *experience) (*resiliency1.Percentile, *core.Status) {
			// rate limiting state
			e, status := exp.getRateLimitingAction(r.handler, r.origin)
			if status.OK() {
				r.state.rateLimit = e.Limit
				r.state.rateBurst = e.Burst
			} else {
				r.handler.Handle(status, "")
			}
			// start ticker
			r.startup()
			return r.percentile, status
		},
		process: func(r *resiliency, percentile *resiliency1.Percentile, observe *observation, exp *experience) ([]timeseries1.Entry, *core.Status) {
			r.handler.AddActivity(r.agentId, "onTick")
			ts, status1 := observe.timeseries(r.handler, r.origin)
			if !status1.OK() || status1.NotFound() {
				return ts, status1
			}
			i, status := resiliencyInference(r, ts, percentile)
			if !status.OK() {
				return ts, status
			}
			status = exp.addInference(r.handler, r.origin, i)
			if !status.OK() {
				return ts, status
			}
			action, status2 := resiliencyAction(r, i)
			if !status2.OK() {
				return ts, status2
			}
			status = exp.addRateLimitingAction(r.handler, r.origin, action)
			return ts, status
		},
		inference: resiliencyInference,
		action:    resiliencyAction,
	}
}()

func resiliencyInference(c *resiliency, entry []timeseries1.Entry, percentile *resiliency1.Percentile) (inference1.Entry, *core.Status) {

	return inference1.Entry{}, core.StatusOK()
}

func resiliencyAction(r *resiliency, entry inference1.Entry) (action1.RateLimiting, *core.Status) {
	return action1.RateLimiting{}, core.StatusOK()
}