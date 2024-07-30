package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
)

// return HTTP status no content if no inference generated
func infer(c *controller, entry []access1.Entry, percentile percentile1.Entry, exp *experience) (inference1.Entry, *core.Status) {

	return inference1.Entry{}, core.StatusOK()
}
