package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
)

type inference struct {
	process func(e []access1.Entry, percentile percentile1.Entry, observe *observation) *core.Status
}

func newInference(handler func(status *core.Status, _ string) *core.Status) *inference {
	if handler == nil {
		handler = func(status *core.Status, _ string) *core.Status {
			return status
		}
	}
	return &inference{
		process: func(e []access1.Entry, percentile percentile1.Entry, observe *observation) *core.Status {
			i, status := infer(e, percentile, observe)
			if !status.OK() {
				handler(status, "")
				return status
			}
			return observe.addInference(i)
		},
	}
}

func infer(entry []access1.Entry, percentile percentile1.Entry, observe *observation) (inference1.Entry, *core.Status) {
	return inference1.Entry{}, core.StatusOK()
}
