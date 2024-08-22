package ingress1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"time"
)

type redirectFunc struct {
	process func(r *redirect, observe *common.Observation, guide *common.Guidance) (completed bool, status *core.Status)
	update  func(r *redirect, guide *common.Guidance, local *guidance, ok bool) *core.Status
}

var redirection = func() *redirectFunc {
	return &redirectFunc{
		// TODO : based on process need to do the following:
		// 1. Process observation and determine if SLO is met
		// 2. If SLO is met, then update percentage
		// 3. How to update percentage
		process: func(r *redirect, observe *common.Observation, guide *common.Guidance) (completed bool, status *core.Status) {
			redirectOrigin := r.origin
			redirectOrigin.Host = r.state.Location
			_, status = observe.IngressTimeseries(r.handler, redirectOrigin)
			if !status.OK() {
				return true, status
			}
			// Need to verify that observation meets the percentile SLO
			// if the observation meets the SLO, then create a new Routing action
			action := resiliency1.RoutingAction{
				EntryId:     r.state.EntryId,
				RouteName:   r.state.RouteName,
				CreatedTS:   time.Time{},
				InferenceId: 0,
				Location:    r.state.Location,
				Percentage:  r.state.Percentage,
			}
			status = guide.AddRoutingAction(r.handler, r.origin, &action)
			if !completed {
				r.updatePercentage()
			}
			return completed, status
		},
		update: func(r *redirect, guide *common.Guidance, local *guidance, ok bool) *core.Status {
			rs := resiliency1.RedirectStatusSucceeded
			if !ok {
				rs = resiliency1.RedirectStatusFailed
			}
			status := local.updateRedirect(r.handler, r.origin, rs)
			if !status.OK() {
				return status
			}
			status = guide.AddRedirectAction(r.handler, r.origin, &resiliency1.RedirectAction{
				EntryId:     r.state.EntryId,
				RouteName:   r.state.RouteName,
				CreatedTS:   time.Time{},
				InferenceId: 0,
				Location:    r.state.Location,
				StatusCode:  r.state.Status,
			})
			return status
		},
	}
}()
