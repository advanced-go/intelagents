package ingress1

import (
	"github.com/advanced-go/experience/action1"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

// return HTTP status no content if no inference generated
func controllerInference(c *controller, entry []timeseries1.Entry, percentile percentile1.Entry) (inference1.Entry, *core.Status) {

	return inference1.Entry{}, core.StatusOK()
}

// return HTTP status no content if no action generated
func controllerAction(c *controller, entry inference1.Entry) ([]action1.Entry, *core.Status) {
	return []action1.Entry{}, core.StatusOK()
}

// return HTTP status no content if no action generated
func controllerReviseTicker(c *controller) {
	//return []action1.Entry{}, core.StatusOK()
}
