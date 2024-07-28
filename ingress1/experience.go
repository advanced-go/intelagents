package ingress1

import (
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

// A nod to Linus Torvalds and plain C
type experience struct {
	reviseTicker func(update func(duration time.Duration), ops *operations)
}

func newExperience(agent messaging.OpsAgent) *experience {
	return &experience{
		reviseTicker: func(update func(duration time.Duration), ops *operations) {

			//c.updateTicker(0)
			update(0)

		},
	}
}
