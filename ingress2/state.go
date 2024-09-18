package ingress2

import (
	"github.com/advanced-go/access/threshold1"
)

type observation struct {
	actual threshold1.Entry
	limit  threshold1.Entry
}

func newObservation(actual, limit threshold1.Entry) *observation {
	o := new(observation)
	o.actual = actual
	o.limit = limit
	return o
}
