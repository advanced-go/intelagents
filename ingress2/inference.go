package ingress2

import (
	"errors"
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/intelagents/common2"
	"github.com/advanced-go/stdlib/core"
)

func runInference(r *resiliency, observe *observation, exp *common2.Experience) (*inference1.Entry, *action1.RateLimiting, *core.Status) {
	if observe == nil || exp == nil {
		return nil, nil, core.NewStatusError(core.StatusInvalidArgument, errors.New("error: input is nil"))
	}
	inf := newInference(observe)
	status := exp.AddInference(r.handler, r.origin, *inf)
	if !status.OK() {
		return inf, nil, status
	}
	// Check to see if inference is not actionable, then return
	act := newAction(inf)
	status = exp.AddRateLimitingAction(r.handler, r.origin, *act)

	return inf, act, status
}

func newInference(observe *observation) *inference1.Entry {
	inf := new(inference1.Entry)
	return inf
}

func newAction(inf *inference1.Entry) *action1.RateLimiting {
	act := action1.NewRateLimiting()
	return act
}
