package ingress1

import (
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

// A nod to Linus Torvalds and plain C
type resiliencyFunc struct {
	startup func(r *resiliency, guide *guidance) *core.Status
	process func(r *resiliency, observe *common.Observation, exp *common.Experience, guide *common.Guidance) ([]timeseries1.Entry, *core.Status)
	//inference func(r *resiliency, entry []timeseries1.Entry) (inference1.Entry, *core.Status)
	//action    func(r *resiliency, entry inference1.Entry) (resiliency1.RateLimitingAction, *core.Status)
}

var (
	action = func(r *resiliency, entry inference1.Entry) (resiliency1.RateLimitingAction, *core.Status) {
		return resiliencyAction(r, entry)
	}
	inference = func(r *resiliency, entry []timeseries1.Entry) (inference1.Entry, *core.Status) {
		return resiliencyInference(r, entry)
	}
	resilience = func() *resiliencyFunc {
		return &resiliencyFunc{
			startup: func(r *resiliency, guide *guidance) *core.Status {
				s, status := guide.resiliencyState(r.handler, r.origin)
				if status.OK() {
					*r.state = *s
				}
				r.startup()
				return status
			},
			process: func(r *resiliency, observe *common.Observation, exp *common.Experience, guide *common.Guidance) ([]timeseries1.Entry, *core.Status) {
				ts, status1 := observe.IngressTimeseries(r.handler, r.origin)
				if !status1.OK() || status1.NotFound() {
					return ts, status1
				}
				i, status := inference(r, ts)
				if !status.OK() {
					return ts, status
				}
				status = exp.AddInference(r.handler, r.origin, i)
				if !status.OK() {
					return ts, status
				}
				action, status2 := action(r, i)
				if !status2.OK() {
					return ts, status2
				}
				status = guide.AddRateLimitingAction(r.handler, r.origin, &action)
				return ts, status
			},
		}
	}()
)

func resiliencyInference(c *resiliency, entry []timeseries1.Entry) (inference1.Entry, *core.Status) {
	return inference1.Entry{}, core.StatusOK()
}

func resiliencyAction(r *resiliency, entry inference1.Entry) (resiliency1.RateLimitingAction, *core.Status) {
	return resiliency1.RateLimitingAction{}, core.StatusOK()
}
