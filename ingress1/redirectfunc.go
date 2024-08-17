package ingress1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/stdlib/core"
)

// A nod to Linus Torvalds and plain C
type redirectFunc struct {
	startup func(r *redirect, guide *guidance) (*resiliency1.IngressRedirectState, *core.Status)
	process func(r *redirect, observe *observation) *core.Status
}

var redirection = func() *redirectFunc {
	return &redirectFunc{
		startup: func(r *redirect, guide *guidance) (*resiliency1.IngressRedirectState, *core.Status) {
			s, status := guide.redirectState(r.handler, r.origin)
			r.startup()
			return s, status
		},
		process: func(r *redirect, observe *observation) *core.Status {
			r.handler.AddActivity(r.agentId, "onProcess()")
			redirectOrigin := r.origin
			redirectOrigin.Host = r.state.Location
			_, status := observe.timeseries(r.handler, redirectOrigin)
			if !status.OK() {
				return status
			}
			// Need to verify that observation meets the percentile SLO
			return core.StatusOK()
		},
	}
}()
