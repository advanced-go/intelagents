package ingress2

import (
	"github.com/advanced-go/events/timeseries1"
	"github.com/advanced-go/intelagents/common1"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

// run - ingress resiliency for the LHC
func runResiliencyLHC(r *resiliency, observe *common1.Events) {
	ticker := messaging.NewTicker(r.duration)
	limit := timeseries1.Threshold{}
	common1.SetPercentileThreshold(r.handler, r.origin, &limit, observe)

	ticker.Start(-1)
	for {
		// observation processing
		select {
		case <-ticker.C():
			actual, status := observe.PercentThresholdQuery(r.handler, r.origin, time.Now().UTC(), time.Now().UTC())
			if status.OK() {
				m := messaging.NewRightChannelMessage("", r.agentId, messaging.ObservationEvent, common1.NewObservation(actual, limit))
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
				r.lhc.Close()
				return
			case messaging.DataChangeEvent:
				if p := common1.GetProfile(r.handler, r.agentId, msg); p != nil {
					ticker.Start(p.ResiliencyDuration(-1))
				}
			default:
				r.handler.Handle(common1.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}
