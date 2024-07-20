package egress1

import (
	"context"
	"github.com/advanced-go/guidance/controller1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	queryControllerDuration = time.Second * 2
)

type guidance struct {
	controllers func(origin core.Origin) ([]controller1.Rowset, *core.Status)
}

func newGuidance() *guidance {
	return &guidance{
		controllers: func(origin core.Origin) ([]controller1.Rowset, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), queryControllerDuration)
			defer cancel()
			return controller1.Query(ctx, origin)
		},
	}
}
