package ingress1

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

type resiliencyFunc struct {
	startup func(r *resiliency, guide *common.Guidance) *core.Status
	process func(r *resiliency, observe *common.Observation, exp *common.Experience) ([]timeseries1.Entry, *core.Status)
}

var (
	newAction = func(r *resiliency, entry inference1.Entry) (action1.RateLimiting, *core.Status) {
		return resiliencyAction(r, entry)
	}
	newInference = func(r *resiliency, entry []timeseries1.Entry) (inference1.Entry, *core.Status) {
		return resiliencyInference(r, entry)
	}
	resilience = func() *resiliencyFunc {
		return &resiliencyFunc{
			startup: func(r *resiliency, guide *common.Guidance) *core.Status {
				s, status := guide.ResiliencyState(r.handler, r.origin)
				if status.OK() {
					r.state = s
				}
				r.startup()
				return status
			},
			process: func(r *resiliency, observe *common.Observation, exp *common.Experience) ([]timeseries1.Entry, *core.Status) {
				ts, status1 := observe.IngressTimeseries(r.handler, r.origin)
				if !status1.OK() || status1.NotFound() {
					return ts, status1
				}
				i, status := newInference(r, ts)
				if !status.OK() {
					return ts, status
				}
				status = exp.AddInference(r.handler, r.origin, i)
				if !status.OK() {
					return ts, status
				}
				action, status2 := newAction(r, i)
				if !status2.OK() {
					return ts, status2
				}
				status = exp.AddRateLimitingAction(r.handler, r.origin, action)
				return ts, status
			},
		}
	}()
)

func resiliencyInference(c *resiliency, entry []timeseries1.Entry) (inference1.Entry, *core.Status) {
	return inference1.Entry{}, core.StatusOK()
}

func resiliencyAction(r *resiliency, entry inference1.Entry) (action1.RateLimiting, *core.Status) {
	return action1.RateLimiting{}, core.StatusOK()
}
