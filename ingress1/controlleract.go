package ingress1

import (
	"github.com/advanced-go/observation/action1"
	"github.com/advanced-go/observation/inference1"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

// return HTTP status no content if no inference generated
func act(entry inference1.Entry, agent messaging.OpsAgent) ([]action1.Entry, *core.Status) {
	return []action1.Entry{}, core.StatusOK()
}
