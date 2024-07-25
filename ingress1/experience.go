package ingress1

import (
	"github.com/advanced-go/stdlib/messaging"
)

// A nod to Linus Torvalds and plain C
type experience struct {
}

func newExperience(agent messaging.OpsAgent) *experience {
	return &experience{}
}
