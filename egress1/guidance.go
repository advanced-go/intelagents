package egress1

import (
	"github.com/advanced-go/stdlib/core"
	"time"
)

const (
	queryControllerDuration = time.Second * 2
)

type guidance struct {
	controllers func(origin core.Origin) *core.Status
}

var (
	localGuidance = func() *guidance {
		return &guidance{}
	}()
)
