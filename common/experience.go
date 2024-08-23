package common

import (
	"context"
	"github.com/advanced-go/experience/inference1"
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	addInferenceDuration = time.Second * 2
)

// Experience - experience functions struct, a nod to Linus Torvalds and plain C
type Experience struct {
	AddInference func(h core.ErrorHandler, origin core.Origin, entry inference1.Entry) *core.Status
}

var Exp = func() *Experience {
	return &Experience{
		AddInference: func(h core.ErrorHandler, origin core.Origin, e inference1.Entry) *core.Status {
			ctx, cancel := context.WithTimeout(context.Background(), addInferenceDuration)
			defer cancel()
			status := inference1.IngressInsert(ctx, nil, e)
			if !status.OK() && !status.NotFound() {
				h.Handle(status, "")
			}
			return status
		},
	}
}()
