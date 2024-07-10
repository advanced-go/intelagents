package guidance

import (
	"context"
	"github.com/advanced-go/guidance/percentile1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	duration = time.Second * 2
)

// TODO : Configure context with deadline

// GetPercentile - resource GET
func GetPercentile[T core.ErrorHandler](origin core.Origin, handle T) percentile1.Entry {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	entry, status := percentile1.Get(ctx, origin)
	handle(status, "")
	return entry
}
