package ingress2

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
)

func runInference(r *resiliency, observe *observation) *inference1.Entry {
	inf := inference1.NewEntry()

	return inf
}

func newAction(inf *inference1.Entry) *action1.RateLimiting {
	act := action1.NewRateLimiting()
	return act
}
