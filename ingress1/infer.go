package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
)

func infer(entry []access1.Entry, percentile percentile1.Entry, observe *observation) (inference1.Entry, *core.Status) {
	return inference1.Entry{}, core.StatusOK()
}
