package ingress1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

type redirectFunc struct {
	startup  func(r *redirect, guide *guidance) (*resiliency1.IngressRedirectState, *core.Status)
	process  func(r *redirect, observe *common.Observation, exp *common.Experience, guide *guidance) *core.Status
	process2 func(r *resiliency, observe *common.Observation, exp *common.Experience) ([]timeseries1.Entry, *core.Status)
}

var redirection = func() *redirectFunc {
	return &redirectFunc{
		startup: func(r *redirect, guide *guidance) (*resiliency1.IngressRedirectState, *core.Status) {
			s, status := guide.redirectState(r.handler, r.origin)
			r.startup()
			return s, status
		},
		// TODO : based on process need to do the following:
		// 1. Update percentage and send action
		// 2. if status fail, then update redirect
		// 3. if status succeed, then update redirect and set redirect action
		// 4. IF done, then message parent and shutdown
		process: func(r *redirect, observe *common.Observation, exp *common.Experience, guide *guidance) *core.Status {
			redirectOrigin := r.origin
			redirectOrigin.Host = r.state.Location
			_, status := observe.IngressTimeseries(r.handler, redirectOrigin)
			if !status.OK() {
				return status
			}
			// Need to verify that observation meets the percentile SLO
			return core.StatusOK()
		},

		process2: func(r *resiliency, observe *common.Observation, exp *common.Experience) ([]timeseries1.Entry, *core.Status) {
			r.handler.AddActivity(r.agentId, "onTick")
			ts, status1 := observe.IngressTimeseries(r.handler, r.origin)
			if !status1.OK() || status1.NotFound() {
				return ts, status1
			}
			i, status := resiliencyInference(r, ts)
			if !status.OK() {
				return ts, status
			}
			status = exp.AddInference(r.handler, r.origin, i)
			if !status.OK() {
				return ts, status
			}
			action, status2 := resiliencyAction(r, i)
			if !status2.OK() {
				return ts, status2
			}
			status = exp.AddRateLimitingAction(r.handler, r.origin, action)
			return ts, status
		},
	}
}()
