package ingress1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

// A nod to Linus Torvalds and plain C
type redirectFunc struct {
	init    func(r *redirect, exp *experience) *core.Status
	process func(r *redirect, percentile resiliency1.Percentile, observe *observation, exp *experience) ([]timeseries1.Entry, *core.Status)
	//inference func(r *redirect, entry []timeseries1.Entry, percentile resiliency1.Percentile) (inference1.Entry, *core.Status)
	//action    func(r *resiliency, entry inference1.Entry) (action1.RateLimiting, *core.Status)
}

var redirection = func() *redirectFunc {
	return &redirectFunc{
		init: func(r *redirect, exp *experience) *core.Status {
			e, status := exp.getRoutingAction(r.handler, r.origin)
			if status.OK() {
				r.state.host = r.origin.Host
				r.state.location = e.Location
				r.state.percent = 0
				return status
			}
			if !status.NotFound() {
				r.handler.Handle(status, "")
			}
			return status
		},
		process: func(r *redirect, percentile resiliency1.Percentile, observe *observation, exp *experience) ([]timeseries1.Entry, *core.Status) {
			r.handler.AddActivity(r.agentId, "onTick")
			/*
				ts, status1 := observe.timeseries(c.handler, c.origin)
				if !status1.OK() || status1.NotFound() {
					return ts, status1
				}
				i, status := exp.processInference(c, ts, percentile)
				if !status.OK() {
					return ts, status
				}
				status = exp.addInference(c.handler, c.origin, i)
				if !status.OK() {
					return ts, status
				}
				actions, status2 := exp.processAction(c, i)
				if !status2.OK() {
					return ts, status2
				}
				status = exp.addAction(c.handler, c.origin, actions)

			*/
			return nil, core.StatusOK()
		},
	}
}()
