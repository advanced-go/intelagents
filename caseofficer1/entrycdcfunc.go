package caseofficer1

import "github.com/advanced-go/stdlib/core"

// A nod to Linus Torvalds and plain C
type entryCDCFunc struct {
	startup func(r *entryCDC, guide *guidance) *core.Status
	//percentile func(r *caseOfficer, curr *resiliency1.Percentile, guide *guidance) (*resiliency1.Percentile, *core.Status)
	//inference func(r *caseOfficer, entry []timeseries1.Entry, percentile *resiliency1.Percentile) (inference1.Entry, *core.Status)
	//action    func(r *caseOfficer, entry inference1.Entry) (action1.RateLimiting, *core.Status)
	//percentile
}

var entry = func() *entryCDCFunc {
	return &entryCDCFunc{
		startup: func(c *entryCDC, guide *guidance) *core.Status {
			//cdc,status := guide.
			return core.StatusOK()
		},

		/*
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

		*/
	}
}()
