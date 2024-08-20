package egress1

import (
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
	"time"
)

const (
	queryControllerDuration = time.Second * 2
)

type guidance struct {
	controllers func(origin core.Origin) *core.Status
}

func newGuidance(agent messaging.OpsAgent) *guidance {
	return &guidance{
		controllers: func(origin core.Origin) *core.Status {
			//ctx, cancel := context.WithTimeout(context.Background(), queryControllerDuration)
			//defer cancel()
			// status := controller1.EgressControllers(ctx, origin)
			///if !status.OK() && !status.NotFound() {
			//	agent.Handle(status, "")
			//	}
			return core.StatusOK()
		},
	}
}
