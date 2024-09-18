package ingress2

import (
	"github.com/advanced-go/access/threshold1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/intelagents/common2"
	"github.com/advanced-go/stdlib/messaging"
)

const (
	thresholdPercent = 95
	thresholdValue   = 3000 // milliseconds
	thresholdMinimum = 0
)

// run - ingress resiliency for the LHC
func runResiliencyLHC(r *resiliency, observe *common2.Access) {
	ticker := messaging.NewTicker(r.duration)
	limit := threshold1.Entry{}
	updateThreshold(r, &limit, observe)

	ticker.Start(-1)
	for {
		// observation processing
		select {
		case <-ticker.C():
			r.handler.AddActivity(r.agentId, "onTick")
			actual, status := observe.Threshold(r.handler, r.origin)
			if status.OK() {
				m := messaging.NewRightChannelMessage("", r.agentId, messaging.ObservationEvent, newObservation(actual[0], limit))
				r.Message(m)
			}
		default:
		}
		// message processing
		select {
		case msg := <-r.lhc.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				ticker.Stop()
				r.rhc.Close()
				return
			case messaging.DataChangeEvent:
				if p := common2.GetProfile(r.handler, r.agentId, msg); p != nil {
					ticker.Start(p.ResiliencyDuration(-1))
				}
			default:
				r.handler.Handle(common.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}

func updateThreshold(r *resiliency, t *threshold1.Entry, observe *common2.Access) {
	if r == nil || t == nil {
		return
	}
	e, status := observe.Threshold(r.handler, r.origin)
	if status.OK() {
		t.Percent = e[0].Percent
		t.Value = e[0].Value
		t.Minimum = e[0].Minimum
	} else {
		t.Percent = thresholdPercent
		t.Value = thresholdValue
		t.Minimum = thresholdMinimum
	}
}
