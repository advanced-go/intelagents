package ingress1

import (
	"context"
	"github.com/advanced-go/observation/action1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	actionInsertDuration = time.Second * 2
)

// A nod to Linus Torvalds and plain C
type action struct {
	process func(entry inference1.Entry, ops *operations) ([]action1.Entry, *core.Status)
	insert  func(origin core.Origin, entry []action1.Entry) *core.Status
}

func newAction(agent messaging.OpsAgent) *action {
	return &action{
		process: func(entry inference1.Entry, ops *operations) ([]action1.Entry, *core.Status) {
			return act(entry, agent, ops)
		},
		insert: func(origin core.Origin, entry []action1.Entry) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), actionInsertDuration)
			defer cancel()
			status := action1.InsertIngress(ctx, origin, entry)
			if !status.OK() && !status.NotFound() {
				agent.Handle(status, "")
			}
			return status
		},
	}
}

// return HTTP status no content if no inference generated
func act(entry inference1.Entry, agent messaging.OpsAgent, ops *operations) ([]action1.Entry, *core.Status) {
	return []action1.Entry{}, core.StatusOK()
}
