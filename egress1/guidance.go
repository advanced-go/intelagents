package egress1

import (
	"context"
	"github.com/advanced-go/guidance/controller1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	queryControllerDuration = time.Second * 2
)

type guidance struct {
	controllers func(origin core.Origin) ([]controller1.Egress, *core.Status)
}

func newGuidance(agent messaging.OpsAgent) *guidance {
	return &guidance{
		controllers: func(origin core.Origin) ([]controller1.Egress, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryControllerDuration)
			defer cancel()
			r, status := controller1.EgressControllers(ctx, origin)
			if !status.OK() && !status.NotFound() {
				agent.Handle(status, "")
			}
			return r, status
		},
	}
}
