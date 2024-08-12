package ingress1

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

// return HTTP status no content if no inference generated
func resiliencyInference(c *resiliency, entry []timeseries1.Entry, percentile resiliency1.Percentile) (inference1.Entry, *core.Status) {

	return inference1.Entry{}, core.StatusOK()
}

// return HTTP status no content if no action generated
func resiliencyAction(r *resiliency, entry inference1.Entry) ([]action1.Entry, *core.Status) {
	return []action1.Entry{}, core.StatusOK()
}

// return HTTP status no content if no action generated
func resiliencyReviseTicker(r *resiliency) {
	//return []action1.Entry{}, core.StatusOK()
}
