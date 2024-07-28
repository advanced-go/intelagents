package ingress1

import (
	"github.com/advanced-go/stdlib/messaging"
)

// A nod to Linus Torvalds and plain C
type experience struct {
	reviseTicker func(c *controller, ops *operations)
}

func newExperience(agent messaging.OpsAgent) *experience {
	return &experience{
		reviseTicker: func(c *controller, ops *operations) {

			//c.updateTicker(nil)
		},
	}
}
