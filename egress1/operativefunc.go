package egress1

import (
	"github.com/advanced-go/guidance/resiliency1"
	"github.com/advanced-go/intelagents/common"
	"github.com/advanced-go/stdlib/core"
	"github.com/advanced-go/stdlib/messaging"
)

type operativeFunc struct {
	startup func(f *fieldOperative, guide *common.Guidance) *core.Status
	update  func(f *fieldOperative, guide *guidance) *core.Status
}

var (
	newAgent = func(origin core.Origin, profile *common.Profile, state *resiliency1.EgressState, handler messaging.OpsAgent) messaging.Agent {
		return newResiliencyAgent(origin, profile, state, handler)
	}

	operative = func() *operativeFunc {
		return &operativeFunc{
			startup: func(f *fieldOperative, guide *common.Guidance) *core.Status {
				s, status := guide.EgressState(f.handler, f.origin)
				if !status.OK() {
					return status
				}
				updateExchange(f, s)
				return status
			},
		}
	}()
)

func updateExchange(f *fieldOperative, entries []resiliency1.EgressState) {
	for i, _ := range entries {
		e := entries[i]
		o := core.Origin{
			Region:     e.Region,
			Zone:       e.Zone,
			SubZone:    "",
			Host:       "",
			InstanceId: "",
			Route:      "",
		}
		agent := newResiliencyAgent(o, f.profile, &e, f.handler)
		f.agents.Register(agent)
		agent.Run()
	}
}
