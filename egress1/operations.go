package egress1

import (
	"github.com/advanced-go/stdlib/messaging"
)

const ()

type operations struct {
	addActivity func(uri string, content any)
}

func newOperations(agent messaging.OpsAgent) *operations {
	return &operations{
		addActivity: func(agentId string, content any) {
			agent.AddActivity(agentId, content)
		},
	}
}
