package guidance1

import (
	"context"
	"github.com/advanced-go/guidance/schedule1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	PercentilePollingDuration = time.Hour * 12
)

func ShouldProcess() bool {
	return true
}

func proc(origin core.Origin, h core.ErrorHandler) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_, status := schedule1.Get(ctx, origin)
	h.Handle(status, "")
	return true
}
