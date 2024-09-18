package common2

import "github.com/advanced-go/events/threshold1"

type Observation struct {
	Actual threshold1.Entry
	Limit  threshold1.Entry
}

func NewObservation(actual, limit threshold1.Entry) *Observation {
	o := new(Observation)
	o.Actual = actual
	o.Limit = limit
	return o
}
