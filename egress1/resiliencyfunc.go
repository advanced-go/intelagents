package egress1

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

type resiliencyFunc struct {
	startup func(r *resiliency, guide *common.Guidance) *core.Status
	process func(r *resiliency, observe *common.Observation, exp *common.Experience, guide *common.Guidance) ([]timeseries1.Entry, *core.Status)
}

var (
	processRateLimitingAction = func(r *resiliency, entry inference1.Entry) (action1.RateLimiting, *core.Status) {
		return resiliencyRateLimitingAction(r, entry)
	}
	processRoutingAction = func(r *resiliency, entry inference1.Entry) (action1.Routing, *core.Status) {
		return resiliencyRoutingAction(r, entry)
	}
	processInference = func(r *resiliency, entry []timeseries1.Entry) (inference1.Entry, *core.Status) {
		return resiliencyInference(r, entry)
	}
	resilience = func() *resiliencyFunc {
		return &resiliencyFunc{
			startup: func(r *resiliency, guide *common.Guidance) *core.Status {
				//s, status := guide.ResiliencyState(r.handler, r.origin)
				//if status.OK() {
				//	*r.state = *s
				//}
				r.startup()
				return core.StatusOK()
			},
			process: func(r *resiliency, observe *common.Observation, exp *common.Experience, guide *common.Guidance) ([]timeseries1.Entry, *core.Status) {
				ts, status1 := observe.IngressTimeseries(r.handler, r.origin)
				if !status1.OK() || status1.NotFound() {
					return ts, status1
				}
				i, status := processInference(r, ts)
				if !status.OK() {
					return ts, status
				}
				status = exp.AddInference(r.handler, r.origin, i)
				if !status.OK() {
					return ts, status
				}
				action, status2 := processRateLimitingAction(r, i)
				if !status2.OK() {
					return ts, status2
				}
				status = exp.AddRateLimitingAction(r.handler, r.origin, action)

				// TODO : add processing of routing action
				return ts, status
			},
		}
	}()
)

func resiliencyInference(c *resiliency, entry []timeseries1.Entry) (inference1.Entry, *core.Status) {
	return inference1.Entry{}, core.StatusOK()
}

func resiliencyRateLimitingAction(r *resiliency, entry inference1.Entry) (action1.RateLimiting, *core.Status) {
	return action1.RateLimiting{}, core.StatusOK()
}

func resiliencyRoutingAction(r *resiliency, entry inference1.Entry) (action1.Routing, *core.Status) {
	return action1.Routing{}, core.StatusOK()
}
