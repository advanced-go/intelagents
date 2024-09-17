package ingress2

import (
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/observation/timeseries1"
	"github.com/advanced-go/stdlib/core"
)

type resiliencyFunc struct {
	startup func(r *resiliency, guide *common.Guidance) *core.Status
	process func(r *resiliency, observe *common.Observation, exp *common.Experience) ([]timeseries1.Entry, *core.Status)
}
