package ingress1

import (
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/observation/access1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

type inference struct {
	process      func(c *controller, e []access1.Entry, percentile percentile1.Entry, exp *experience, ops *operations) (inference1.Entry, *core.Status)
	updateTicker func(c *controller, exp *experience)
}

func newInference(agent messaging.OpsAgent) *inference {
	return &inference{
		process: func(c *controller, entry []access1.Entry, percentile percentile1.Entry, exp *experience, ops *operations) (inference1.Entry, *core.Status) {
			return infer(c, entry, percentile, exp, agent, ops)
		},
		updateTicker: func(c *controller, exp *experience) {
			c.updateTicker(exp)
		},
	}
}

// return HTTP status no content if no inference generated
func infer(c *controller, entry []access1.Entry, percentile percentile1.Entry, exp *experience, agent messaging.OpsAgent, ops *operations) (inference1.Entry, *core.Status) {
	return inference1.Entry{}, core.StatusOK()
}

/*
func(e []access1.Entry, exp *experience, percentile1.Entry, ) (inference1.Entry,*core.Status) {
	i, status := infer(e, percentile, observe)
	if !status.OK() {
		agent.Handle(status, "")
		return status
	}
	return observe.addInference(i)
},

*/
