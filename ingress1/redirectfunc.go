package ingress1

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"time"
)

type redirectFunc struct {
	process func(r *redirect, observe *common.Observation, exp *common.Experience) (completed bool, status *core.Status)
	update  func(r *redirect, exp *common.Experience, guide *common.Guidance, ok bool) *core.Status
}

var redirection = func() *redirectFunc {
	return &redirectFunc{
		// TODO : based on process need to do the following:
		// 1. Process observation and determine if SLO is met
		// 2. If SLO is met, then update percentage
		// 3. How to update percentage
		process: func(r *redirect, observe *common.Observation, exp *common.Experience) (completed bool, status *core.Status) {
			redirectOrigin := r.origin
			redirectOrigin.Host = r.state.Location
			_, status = observe.IngressTimeseries(r.handler, redirectOrigin)
			if !status.OK() {
				return true, status
			}
			// Need to verify that observation meets the percentile SLO
			// if the observation meets the SLO, then create a new Routing action
			action := action1.Routing{
				EntryId:     r.state.EntryId,
				RouteName:   r.state.Route,
				CreatedTS:   time.Time{},
				InferenceId: 0,
				Location:    r.state.Location,
				Percentage:  r.state.Percentage,
			}
			status = exp.AddRoutingAction(r.handler, r.origin, action)
			if !completed {
				r.updatePercentage()
			}
			return completed, status
		},
		update: func(r *redirect, exp *common.Experience, guide *common.Guidance, ok bool) *core.Status {
			rs := resiliency1.RedirectStatusSucceeded
			if !ok {
				rs = resiliency1.RedirectStatusFailed
			}
			status := guide.UpdateRedirect(r.handler, r.origin, rs)
			if !status.OK() {
				return status
			}
			status = exp.AddRedirectAction(r.handler, r.origin, action1.Redirect{
				EntryId:     r.state.EntryId,
				RouteName:   r.state.Route,
				CreatedTS:   time.Time{},
				InferenceId: 0,
				Location:    r.state.Location,
				StatusCode:  r.state.Status,
			})
			return status
		},
	}
}()
