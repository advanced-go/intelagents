package ingress1

import "github.com/advanced-go/stdlib/messaging"

// A nod to Linus Torvalds and plain C
type action struct {
}

func newAction(agent messaging.OpsAgent) *action {
	return &action{}
}
