package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

type inference struct {
	process func(e []access1.Entry, percentile percentile1.Entry, observe *observation) *core.Status
}

func newInference(agent messaging.OpsAgent) *inference {
	return &inference{
		process: func(e []access1.Entry, percentile percentile1.Entry, observe *observation) *core.Status {
			i, status := infer(e, percentile, observe)
			if !status.OK() {
				agent.Handle(status, "")
				return status
			}
			return observe.addInference(i)
		},
	}
}

func infer(entry []access1.Entry, percentile percentile1.Entry, observe *observation) (inference1.Entry, *core.Status) {
	return inference1.Entry{}, core.StatusOK()
}
