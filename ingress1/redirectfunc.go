package ingress1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"time"
)

type redirectFunc struct {
	startup func(r *redirect, guide *guidance) (*resiliency1.IngressRedirectState, *core.Status)
	process func(r *redirect, observe *common.Observation, guide *common.Guidance) (completed bool, status *core.Status)
}

var redirection = func() *redirectFunc {
	return &redirectFunc{
		startup: func(r *redirect, guide *guidance) (*resiliency1.IngressRedirectState, *core.Status) {
			s, status := guide.redirectState(r.handler, r.origin)
			r.startup()
			return s, status
		},
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
	}
}()
