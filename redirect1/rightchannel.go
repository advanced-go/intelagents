package redirect1

import (
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/messaging"
)

func runRedirectRHC(r *redirect, fn *redirectFunc, observe *common.Observation, exp *common.Experience, guide *common.Guidance) {
	r.startup()
	for {
		select {
		case <-r.ticker.C():
			r.handler.AddActivity(r.agentId, "onTick()")
			completed, status := fn.process(r, observe, exp)
			if completed {
				fn.update(r, exp, guide, status.OK())
				func (r *redirect) updatePercentage() {
					switch r.state.Percentage {
					case 0:
						r.state.Percentage = 10
					case 10:
						r.state.Percentage = 20
					case 20:
						r.handler.Message(messaging.NewControlMessage("", r.agentId, RedirectCompletedEvent))
						r.shutdown()
						return
					}
		case msg := <-r.rhc.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.shutdown()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			default:
			}
		default:
		}
	}
}

